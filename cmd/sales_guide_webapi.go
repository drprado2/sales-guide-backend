package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/drprado2/react-redux-typescript/configs"
	"github.com/drprado2/react-redux-typescript/internal/adapters"
	"github.com/drprado2/react-redux-typescript/internal/domain"
	"github.com/drprado2/react-redux-typescript/internal/domain/entities"
	"github.com/drprado2/react-redux-typescript/internal/domain/valueobjects"
	apptracer2 "github.com/drprado2/react-redux-typescript/pkg/apptracer"
	http_server2 "github.com/drprado2/react-redux-typescript/pkg/httpserver"
	logs2 "github.com/drprado2/react-redux-typescript/pkg/logs"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	zipkintracer "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	ctx := context.Background()
	logs2.Setup()

	configs := configs.Get()

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?connect_timeout=15&application_name=poker-simulator",
		configs.DbUser,
		configs.DbPass,
		configs.DbHost,
		configs.DbPort,
		configs.DbName)
	logs2.Logger(ctx).Info("Connecting postgres DB")
	dbpool, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		logs2.Logger(ctx).Fatalf("error creating DB connection, %v", err)
		os.Exit(1)
	}
	logs2.Logger(ctx).Info("DB connected successfully")
	defer dbpool.Close()

	endpoint, err := zipkin.NewEndpoint("sales-guide", fmt.Sprintf("0.0.0.0:%v", configs.ServerPort))
	if err != nil {
		logs2.Logger(ctx).Fatalf("unable to create local endpoint: %+v\n", err)
	}

	httpReporter := zipkinhttp.NewReporter(configs.ZipkinUrl, zipkinhttp.BatchInterval(time.Second*3))
	defer httpReporter.Close()

	tracer, err := zipkin.NewTracer(
		httpReporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithTraceID128Bit(true),
		zipkin.WithNoopTracer(true))
	if err != nil {
		logs2.Logger(ctx).Fatalf("unable to create apptracer: %+v\n", err)
	}

	tracerService := apptracer2.TracerService{
		Endpoint: endpoint,
		Tracer:   tracer,
	}

	opentracing.SetGlobalTracer(zipkintracer.Wrap(tracer))

	companyRepository := adapters.NewCompanySqlRepository(dbpool, tracerService)
	logger := new(adapters.LogrusLogger)
	companyHttpAdapter := adapters.NewCompanyHttpAdapter(logger)
	cnpjValidator := new(adapters.PaemureBrDocCnpjValidator)
	colorService := new(adapters.CssColorParserService)

	entities.Setup(companyRepository)
	valueobjects.Setup(colorService, cnpjValidator)
	domain.Setup(
		logger,
		tracerService,
	)

	server := http_server2.NewServer(logger, tracer, endpoint).
		WithRoutes(companyHttpAdapter.RegisterRouteHandlers)

	go func() {
		server.Start()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	server.Shutdown(ctx)
	logs2.Logger(ctx).Info("shutting down")
	os.Exit(0)
}
