package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/drprado2/react-redux-typescript/configs"
	"github.com/drprado2/react-redux-typescript/internal/apptracer"
	"github.com/drprado2/react-redux-typescript/internal/http_server"
	"github.com/drprado2/react-redux-typescript/internal/logs"
	requesthandlers "github.com/drprado2/react-redux-typescript/internal/request_handlers"
	"github.com/drprado2/react-redux-typescript/internal/services/players"
	"github.com/drprado2/react-redux-typescript/internal/storage/repositories"
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

	configs := &configs.Envs{}

	logs.Init(configs)

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?connect_timeout=15&application_name=poker-simulator",
		configs.GetDbUser(),
		configs.GetDbPassword(),
		configs.GetDbHost(),
		configs.GetDbPort(),
		configs.GetDbName())
	logs.Logger(ctx).Info("Connecting postgres DB")
	dbpool, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		logs.Logger(ctx).Fatalf("error creating DB connection, %v", err)
		os.Exit(1)
	}
	logs.Logger(ctx).Info("DB connected successfully")
	defer dbpool.Close()

	endpoint, err := zipkin.NewEndpoint("poker-simulator", fmt.Sprintf("0.0.0.0:%v", configs.GetServerPort()))
	if err != nil {
		logs.Logger(ctx).Fatalf("unable to create local endpoint: %+v\n", err)
	}

	httpReporter := zipkinhttp.NewReporter(configs.GetZipkinUrl(), zipkinhttp.BatchInterval(time.Second*3))
	defer httpReporter.Close()

	tracer, err := zipkin.NewTracer(httpReporter, zipkin.WithLocalEndpoint(endpoint), zipkin.WithTraceID128Bit(true))
	if err != nil {
		logs.Logger(ctx).Fatalf("unable to create apptracer: %+v\n", err)
	}

	tracerService := apptracer.TracerService{
		Endpoint: endpoint,
		Tracer:   tracer,
	}

	opentracing.SetGlobalTracer(zipkintracer.Wrap(tracer))

	playerRepository := repositories.PlayerRepository{
		DbPool: dbpool,
		Tracer: tracerService,
	}
	userService := &players.UserService{
		PlayerRepository: playerRepository,
		Tracer: tracerService,
	}
	userHandler := &requesthandlers.UserHandler{
		UserService: userService,
		Tracer: tracerService,
	}

	server := &http_server.Server{
		Envs:           configs,
		UserHandlers:   userHandler,
		Tracer:         tracer,
		ZipkinEndpoint: endpoint,
	}

	go func() {
		server.Start()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	server.Shutdown(ctx)
	logs.Logger(ctx).Info("shutting down")
	os.Exit(0)
}
