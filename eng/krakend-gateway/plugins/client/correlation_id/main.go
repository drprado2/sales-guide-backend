package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func init() {
	fmt.Println("correlation_id plugin is loaded!")
}

func main() {}

// HandlerRegisterer is the name of the symbol krakend looks up to try and register plugins
var HandlerRegisterer registrable = registrable("correlation_id")

type registrable string

const outputHeaderName = "X-Cid"

func (r registrable) RegisterHandlers(f func(
	name string,
	handler func(
		context.Context,
		map[string]interface{},
		http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (r registrable) registerHandlers(ctx context.Context, extra map[string]interface{}, handler http.Handler) (http.Handler, error) {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cid := uuid.NewString()
		r2 := new(http.Request)
		*r2 = *r
		r2.Header.Set(outputHeaderName, cid)
		handler.ServeHTTP(w, r2)
	}), nil
}
