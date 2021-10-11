package middlewares

import (
	"context"
	"fmt"
	"github.com/drprado2/react-redux-typescript/internal/models"
	"github.com/drprado2/react-redux-typescript/pkg/ctxvals"
	logs2 "github.com/drprado2/react-redux-typescript/pkg/logs"
	"github.com/google/uuid"
	"github.com/openzipkin/zipkin-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

const (
	DefaultTimezone = "America/Sao_Paulo"
	DefaultTimeOff  = -3
)

var (
	LocationsCache     = make(map[string]*time.Location)
	DefaultLocation, _ = time.LoadLocation(DefaultTimezone)
	locationMutex      = &sync.Mutex{}
)

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logs2.Logger(r.Context()).Errorf("Panic occurs in path %v, error: %v, stacktrace: %s", r.RequestURI, err, string(debug.Stack()))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func RequestLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCtx := context.WithValue(r.Context(), "httpMethod", r.Method)
		reqCtx = context.WithValue(reqCtx, "httpPath", r.RequestURI)

		logs2.Logger(reqCtx).Info(r.Context(), fmt.Sprintf("Handling request, method: %v, path: %v", r.Method, r.RequestURI), nil)

		writter := &models.StatusRecorder{
			w, http.StatusOK,
		}
		next.ServeHTTP(writter, r.WithContext(reqCtx))

		if writter.Status >= 500 {
			logs2.Logger(r.Context()).WithFields(logrus.Fields{"httpStatusCode": writter.Status, "requestSuccess": false}).Errorf("request fineshed with errors, status code: %v", writter.Status)
		} else if writter.Status >= 400 {
			logs2.Logger(r.Context()).WithFields(logrus.Fields{"httpStatusCode": writter.Status, "requestSuccess": false}).Warnf("request fineshed with errors, status code: %v", writter.Status)
		} else {
			logs2.Logger(r.Context()).WithFields(logrus.Fields{"httpStatusCode": writter.Status, "requestSuccess": true}).Infof("request fineshed with success, status code: %v", writter.Status)
		}
	})
}

func CidMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cid := r.Header.Get("x-cid")
		if cid == "" {
			cid = uuid.NewString()
		}
		reqCtx := context.WithValue(r.Context(), "cid", cid)
		next.ServeHTTP(w, r.WithContext(reqCtx))
	})
}

func SpanMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := zipkin.SpanFromContext(r.Context())
		span.SetName(fmt.Sprintf("%s::%s", r.Method, r.RequestURI))
		span.Tag("cid", r.Context().Value("cid").(string))
		defer span.Finish()

		next.ServeHTTP(w, r.WithContext(zipkin.NewContext(r.Context(), span)))
	})
}

func LocationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timezone := r.Header.Get("x-timezone")
		timeoffset := r.Header.Get("x-timezone-offset")
		iTimeoffset := DefaultTimeOff
		if timezone == "" {
			timezone = DefaultTimezone
		}
		if timeoffset != "" {
			if v, err := strconv.Atoi(timeoffset); err == nil {
				iTimeoffset = v
			}
		}
		reqCtx := ctxvals.WithTimezone(r.Context(), timezone)
		reqCtx = ctxvals.WithTimeOffset(reqCtx, iTimeoffset)

		location, ok := LocationsCache[timezone]
		if !ok {
			loc, err := time.LoadLocation(timezone)
			if err != nil {
				logs2.Logger(r.Context()).Warnf("invalid location err=%v", err)
				loc = DefaultLocation
			}
			newCache := make(map[string]*time.Location)
			for k, v := range LocationsCache {
				newCache[k] = v
			}
			newCache[timezone] = loc
			locationMutex.Lock()
			LocationsCache = newCache
			locationMutex.Unlock()
			location = loc
		}
		reqCtx = ctxvals.WithLocation(reqCtx, location)
		next.ServeHTTP(w, r.WithContext(reqCtx))
	})
}
