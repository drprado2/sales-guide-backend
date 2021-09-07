package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/drprado2/react-redux-typescript/configs"
	"github.com/drprado2/react-redux-typescript/internal/adapters/http/company"
	"github.com/drprado2/react-redux-typescript/internal/adapters/observability"
	"github.com/drprado2/react-redux-typescript/internal/adapters/repository"
	"github.com/drprado2/react-redux-typescript/internal/adapters/utils"
	"github.com/drprado2/react-redux-typescript/internal/domain/usecases"
	utils2 "github.com/drprado2/react-redux-typescript/internal/utils"
	http_server2 "github.com/drprado2/react-redux-typescript/pkg/httpserver"
	logs2 "github.com/drprado2/react-redux-typescript/pkg/logs"
	"github.com/jackc/pgx/v4/pgxpool"
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

	tracerService, endpoint, err := observability.NewTracer(ctx)
	if err != nil {
		logs2.Logger(ctx).Fatalf("error creating tracer service, %v", err)
		os.Exit(1)
	}

	companyRepository := repository.NewCompanySqlRepository(dbpool, tracerService)
	cnpjValidator := new(utils.PaemureBrDocCnpjValidator)
	colorService := new(utils.CssColorParserService)

	usecases.Setup(companyRepository, tracerService)
	utils2.Setup(colorService, cnpjValidator)

	server := http_server2.NewServer(tracerService.Tracer, endpoint).
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
