package httpserver

import (
	"context"
	"fmt"
	"github.com/drprado2/react-redux-typescript/configs"
	"github.com/drprado2/react-redux-typescript/internal/domain"
	middlewares2 "github.com/drprado2/react-redux-typescript/pkg/httpserver/middlewares"
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
	server         *http.Server
	tracer         *zipkin.Tracer
	zipkinEndpoint *zipkinmodel.Endpoint
	logger         domain.Logger
	mainCtx        context.Context
	cancelCtx      context.CancelFunc
	router         *mux.Router
}

func NewServer(logger domain.Logger, tracer *zipkin.Tracer, zipkinEndpoint *zipkinmodel.Endpoint) *Server {
	server := new(Server)
	server.logger = logger
	server.tracer = tracer
	server.zipkinEndpoint = zipkinEndpoint

	r := mux.NewRouter().StrictSlash(true)
	ctx, cancel := context.WithCancel(context.Background())
	server.mainCtx = ctx
	server.cancelCtx = cancel

	server.router = r.PathPrefix("/api").Subrouter()

	return server
}

func (s *Server) WithRoutes(handler func(router *mux.Router)) *Server {
	s.logger.Infof(s.mainCtx, "registering request handlers")

	handler(s.router)
	return s
}

func (s *Server) Start() {
	envs := configs.Get()
	s.logger.Infof(s.mainCtx, "Starting http server at port %v, with env %v", envs.ServerPort, envs.ServerEnvironment)

	http.Handle("/", s.router)

	s.registerMiddlewares(s.router)

	if envs.ServerEnvironment == configs.DeveloperEnvironment {
		go func() {
			s.logger.Infof(s.mainCtx, "starting pgprof server")
			http.DefaultServeMux.Handle("/debug/fgprof", fgprof.Handler())
			if err := http.ListenAndServe(":6060", nil); err != nil {
				s.logger.Errorf(s.mainCtx, err.Error())
			}
		}()
	}

	muxWithMiddlewares := http.TimeoutHandler(s.router, time.Duration(envs.ServerEndpointTimeout)*time.Second, "timeout occurred!")

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%v", envs.ServerPort),
		Handler:      muxWithMiddlewares,
		WriteTimeout: time.Duration(envs.ServerEndpointTimeout) * time.Second,
		ReadTimeout:  time.Duration(envs.ServerEndpointTimeout) * time.Second,
	}

	if err := s.server.ListenAndServe(); err != nil {
		s.logger.Fatalf(s.mainCtx, "Fail starting http server, %v", err)
		panic(err)
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	logs2.Logger(ctx).Infof("shutting down http server")
	s.server.Shutdown(ctx)
}

func (s *Server) registerMiddlewares(router *mux.Router) {
	router.Use(middlewares2.PanicMiddleware)
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(middlewares2.CidMiddleware)
	router.Use(middlewares2.SpanMiddleware(s.tracer, s.zipkinEndpoint))
	router.Use(middlewares2.RequestLogMiddleware)
}
