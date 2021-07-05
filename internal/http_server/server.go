package http_server

import (
	"context"
	"fmt"
	"github.com/drprado2/react-redux-typescript/configs"
	"github.com/drprado2/react-redux-typescript/internal/http_server/middlewares"
	requesthandlers "github.com/drprado2/react-redux-typescript/internal/request_handlers"
	logs2 "github.com/drprado2/react-redux-typescript/pkg/logs"
	"github.com/felixge/fgprof"
	"github.com/gorilla/mux"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinmodel "github.com/openzipkin/zipkin-go/model"
	"net/http"
	_ "net/http/pprof"
	"time"
)

const (
	Post   = "POST"
	Put    = "PUT"
	Delete = "DELETE"
	Get    = "GET"
)

type Server struct {
	Server         *http.Server
	Envs           configs.EnvsInterface
	UserHandlers   requesthandlers.UserHandlerInterface
	Tracer         *zipkin.Tracer
	ZipkinEndpoint *zipkinmodel.Endpoint
}

func (s *Server) Start() {
	r := mux.NewRouter().StrictSlash(true)
	ctx := context.Background()
	logs2.Logger(ctx).Infof("Starting http Server at port %v, with env %v", s.Envs.GetServerPort(), s.Envs.GetEnvironment())

	logs2.Logger(ctx).Infof("registering middlewares")
	s.registerMiddlewares(r)
	logs2.Logger(ctx).Infof("registering request handlers")
	s.registerHandlers(r)

	if s.Envs.GetEnvironment() == configs.DeveloperEnvironment {
		go func() {
			logs2.Logger(ctx).Infof("starting pgprof server")
			http.DefaultServeMux.Handle("/debug/fgprof", fgprof.Handler())
			if err := http.ListenAndServe(":6060", nil); err != nil {
				logs2.Logger(ctx).WithError(err).Error(err.Error())
			}
		}()
	}

	var muxWithMiddlewares http.Handler
	if s.Envs.GetEnvironment() != configs.DeveloperEnvironment {
		muxWithMiddlewares = http.TimeoutHandler(r, time.Minute*3, "timeout occurred!")
	} else {
		muxWithMiddlewares = http.TimeoutHandler(r, time.Second*15, "timeout occurred!")
	}

	s.Server = &http.Server{
		Addr:    fmt.Sprintf(":%v", s.Envs.GetServerPort()),
		Handler: muxWithMiddlewares,
	}

	if err := s.Server.ListenAndServe(); err != nil {
		logs2.Logger(ctx).WithError(err).Fatalf("Fail starting http Server, %v", err)
		panic(err)
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	logs2.Logger(ctx).Infof("shutting down http Server")
	s.Server.Shutdown(ctx)
}

func (s *Server) registerMiddlewares(router *mux.Router) {
	router.Use(middlewares.PanicMiddleware)
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(middlewares.CidMiddleware)
	router.Use(middlewares.SpanMiddleware(s.Tracer, s.ZipkinEndpoint))
	router.Use(middlewares.RequestLogMiddleware)
}

func (s *Server) registerHandlers(router *mux.Router) {
	apiRouters := router.PathPrefix("/api").Subrouter()

	apiRouters.
		Path("v1/players").
		HandlerFunc(s.UserHandlers.Create).
		Name("create user").
		Methods(Post)
	apiRouters.
		Path("v1/players/{id:[0-9a-fA-F]{8}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{12}}").
		HandlerFunc(s.UserHandlers.Update).
		Name("update user").
		Methods(Put)
	apiRouters.
		Path("v1/players/{id:[0-9a-fA-F]{8}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{12}}").
		HandlerFunc(s.UserHandlers.Delete).
		Name("delete user").
		Methods(Delete)
	apiRouters.
		Path("v1/players/{id:[0-9a-fA-F]{8}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{12}}").
		HandlerFunc(s.UserHandlers.GetById).
		Name("get user by id").
		Methods(Get)
	apiRouters.
		HandleFunc("/v1/players", s.UserHandlers.Get).
		Name("get user paged").
		Methods(Get)
}
