package http

import (
	"context"
	"fmt"
	"github.com/drprado2/sales-guide/configs"
	"github.com/drprado2/sales-guide/internal/domain"
	logs2 "github.com/drprado2/sales-guide/pkg/instrumentation/logs"
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

func NewServer(tracer *zipkin.Tracer, serviceManager *domain.ServiceManager) *Server {
	server := new(Server)

	server.zipkinMiddleware = zipkinhttp.NewServerMiddleware(tracer, zipkinhttp.TagResponseSize(true))

	r := mux.NewRouter().StrictSlash(true)
	ctx, cancel := context.WithCancel(context.Background())
	server.mainCtx = ctx
	server.cancelCtx = cancel

	server.router = r.PathPrefix("/api").Subrouter()

	companyController := NewCompanyController(serviceManager)
	userController := NewUserController(serviceManager)

	return server.
		WithRoutes(companyController.RegisterRouteHandlers).
		WithRoutes(userController.RegisterRouteHandlers)
}

func (s *Server) WithRoutes(handler func(router *mux.Router)) *Server {
	logs2.Logger(s.mainCtx).Info("registering request handlers")

	handler(s.router)
	return s
}

func (s *Server) Start() error {
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
		logs2.Logger(s.mainCtx).WithError(err).Error("Fail starting http server")
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) {
	logs2.Logger(ctx).Infof("shutting down http server")
	s.server.Shutdown(ctx)
}

func (s *Server) registerMiddlewares(router *mux.Router) {
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(PanicMiddleware)
	router.Use(CidMiddleware)
	router.Use(LocationMiddleware)
	router.Use(UserMiddleware)
	router.Use(s.zipkinMiddleware)
	router.Use(SpanMiddleware)
	router.Use(RequestLogMiddleware)
}
