package middlewares

import (
	"context"
	"fmt"
	"github.com/drprado2/react-redux-typescript/internal/models"
	logs2 "github.com/drprado2/react-redux-typescript/pkg/logs"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	zipkinmodel "github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logs2.Logger(r.Context()).Fatalf("Panic occurs in path %v, error: %v", r.RequestURI, err)
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

		if writter.Status >= 400 {
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

func SpanMiddleware(tracer *zipkin.Tracer, zipkinEndpoint *zipkinmodel.Endpoint) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sc := tracer.Extract(b3.ExtractHTTP(r))

			sp := tracer.StartSpan(
				r.RequestURI,
				zipkin.Kind(model.Server),
				zipkin.Parent(sc),
				zipkin.RemoteEndpoint(zipkinEndpoint),
			)

			sp.Tag("cid", r.Context().Value("cid").(string))

			ctx := zipkin.NewContext(r.Context(), sp)

			zipkin.TagHTTPMethod.Set(sp, r.Method)
			zipkin.TagHTTPPath.Set(sp, r.URL.Path)
			if r.ContentLength > 0 {
				zipkin.TagHTTPRequestSize.Set(sp, strconv.FormatInt(r.ContentLength, 10))
			}

			defer sp.Finish()

			next.ServeHTTP(w, r.WithContext(ctx))
		})

	}
}
