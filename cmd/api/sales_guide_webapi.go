package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/drprado2/sales-guide/configs"
	"github.com/drprado2/sales-guide/internal/adapters/http"
	"github.com/drprado2/sales-guide/internal/adapters/repository"
	"github.com/drprado2/sales-guide/internal/adapters/validations"
	"github.com/drprado2/sales-guide/internal/domain"

	"github.com/drprado2/sales-guide/pkg/instrumentation/apptracer"
	"github.com/drprado2/sales-guide/pkg/instrumentation/logs"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinreporter "github.com/openzipkin/zipkin-go/reporter/http"
	"gopkg.in/auth0.v5/management"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	ctx, cancelFunc := context.WithCancel(context.Background())
	logs.Setup()

	envconfigs := configs.Get()

	dbpool, err := repository.CreateConnPool(ctx, envconfigs)
	if err != nil {
		panic(err)
	}
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

	tracerService := &apptracer.TracerService{
		Client: zipkinClient,
		Tracer: tracer,
	}

	companyRepository := repository.NewCompanySqlRepository(dbpool, tracerService)
	userRepository := repository.NewUserSqlRepository(dbpool, tracerService)
	valid, trans := validations.CreateValidatorService()
	auth0Manager, err := management.New(configs.Get().Auth0Domain, management.WithClientCredentials(configs.Get().Auth0ClientID, configs.Get().Auth0ClientSecret))
	if err != nil {
		log.Fatalf("unable to create auth0 manager: %+v\n", err)
	}

	serviceManager := domain.CreateServiceManager(companyRepository, userRepository, tracerService, valid, trans, auth0Manager)

	server := http.NewServer(tracer, serviceManager)

	errCh := make(chan error, 1)
	go func() {
		if err := server.Start(); err != nil {
			errCh <- err
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case <-c:
	case err := <-errCh:
		panic(err)
	}

	ctxShutdown, cancel := context.WithTimeout(ctx, wait)
	cancelFunc()
	defer cancel()
	server.Shutdown(ctxShutdown)
	logs.Logger(ctx).Info("shutting down")
	os.Exit(0)
}
