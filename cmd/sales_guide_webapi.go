package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/drprado2/react-redux-typescript/configs"
	"github.com/drprado2/react-redux-typescript/internal/adapters/http/company"
	"github.com/drprado2/react-redux-typescript/internal/adapters/repository"
	"github.com/drprado2/react-redux-typescript/internal/adapters/utils"
	"github.com/drprado2/react-redux-typescript/internal/domain/usecases"
	utils2 "github.com/drprado2/react-redux-typescript/internal/utils"
	apptracer2 "github.com/drprado2/react-redux-typescript/pkg/apptracer"
	http_server2 "github.com/drprado2/react-redux-typescript/pkg/httpserver"
	logs2 "github.com/drprado2/react-redux-typescript/pkg/logs"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	pt_translations "github.com/go-playground/validator/v10/translations/pt_BR"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinreporter "github.com/openzipkin/zipkin-go/reporter/http"
	"log"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"time"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	ctx := context.Background()
	logs2.Setup()

	envconfigs := configs.Get()

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?connect_timeout=15&application_name=poker-simulator",
		envconfigs.DbUser,
		envconfigs.DbPass,
		envconfigs.DbHost,
		envconfigs.DbPort,
		envconfigs.DbName)
	logs2.Logger(ctx).Info("Connecting postgres DB")
	dbpool, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		logs2.Logger(ctx).Fatalf("error creating DB connection, %v", err)
		os.Exit(1)
	}
	logs2.Logger(ctx).Info("DB connected successfully")
	defer dbpool.Close()

	endpoint, err := zipkin.NewEndpoint(envconfigs.ServiceName, fmt.Sprintf("%s:%v", envconfigs.ServerHost, envconfigs.ServerPort))
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	reporter := zipkinreporter.NewReporter(envconfigs.ZipkinReportUrl)

	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	zipkinClient, err := zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		log.Fatalf("unable to create client: %+v\n", err)
	}

	tracerService := &apptracer2.TracerService{
		Client: zipkinClient,
		Tracer: tracer,
	}

	companyRepository := repository.NewCompanySqlRepository(dbpool, tracerService)
	cnpjValidator := new(utils.PaemureBrDocCnpjValidator)
	colorService := new(utils.CssColorParserService)
	valid, trans := createValidatorService()

	usecases.Setup(companyRepository, tracerService, valid, trans)
	utils2.Setup(colorService, cnpjValidator)

	server := http_server2.NewServer(tracer).
		WithRoutes(company.RegisterRouteHandlers)

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

func createValidatorService() (*validator.Validate, ut.Translator) {
	ptbr := pt_BR.New()
	uni := ut.New(ptbr, ptbr)
	trans, _ := uni.GetTranslator("pt_BR")
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("name"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	if err := pt_translations.RegisterDefaultTranslations(v, trans); err != nil {
		panic(err)
	}
	return v, trans
}
