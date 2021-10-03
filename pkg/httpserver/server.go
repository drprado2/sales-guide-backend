package httpserver

import (
	"context"
	"fmt"
	"github.com/drprado2/react-redux-typescript/configs"
	middlewares2 "github.com/drprado2/react-redux-typescript/pkg/httpserver/middlewares"
	logs2 "github.com/drprado2/react-redux-typescript/pkg/logs"
	"github.com/felixge/fgprof"
	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
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
	server           *http.Server
	zipkinMiddleware func(http.Handler) http.Handler
	mainCtx          context.Context
	cancelCtx        context.CancelFunc
	router           *mux.Router
}

func NewServer(tracer *zipkin.Tracer) *Server {
	server := new(Server)

	server.zipkinMiddleware = zipkinhttp.NewServerMiddleware(tracer, zipkinhttp.TagResponseSize(true))

	r := mux.NewRouter().StrictSlash(true)
	ctx, cancel := context.WithCancel(context.Background())
	server.mainCtx = ctx
	server.cancelCtx = cancel

	server.router = r.PathPrefix("/api").Subrouter()

	return server
}

func (s *Server) WithRoutes(handler func(router *mux.Router)) *Server {
	logs2.Logger(s.mainCtx).Info("registering request handlers")

	handler(s.router)
	return s
}

func (s *Server) Start() {
	envs := configs.Get()
	http.Handle("/", s.router)

	s.registerMiddlewares(s.router)

	if envs.ServerEnvironment == configs.DeveloperEnvironment {
		go func() {
			logs2.Logger(s.mainCtx).Info("starting pgprof server")
			http.DefaultServeMux.Handle("/debug/fgprof", fgprof.Handler())
			if err := http.ListenAndServe(":6060", nil); err != nil {
				logs2.Logger(s.mainCtx).WithError(err).Error("error strarting pgprof")
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

	logs2.Logger(s.mainCtx).Infof("http server running at port %v, with env %v", envs.ServerPort, envs.ServerEnvironment)
	if err := s.server.ListenAndServe(); err != nil {
		logs2.Logger(s.mainCtx).WithError(err).Fatal("Fail starting http server")
		panic(err)
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	logs2.Logger(ctx).Infof("shutting down http server")
	s.server.Shutdown(ctx)
}

func (s *Server) registerMiddlewares(router *mux.Router) {
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(middlewares2.PanicMiddleware)
	router.Use(middlewares2.CidMiddleware)
	router.Use(s.zipkinMiddleware)
	router.Use(middlewares2.SpanMiddleware)
	router.Use(middlewares2.RequestLogMiddleware)
}
