package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

func init() {
	fmt.Println("api_key_authentication plugin is loaded!")
}

func main() {}

// HandlerRegisterer is the name of the symbol krakend looks up to try and register plugins
var HandlerRegisterer registrable = registrable("api-key-auth")

type registrable string

const apiKeyHeader = "X-Api-Key"

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
	configs, okConfigs := extra["api-key-auth"].(map[string]interface{})
	if !okConfigs {
		panic(errors.New("wrong configs of api key plugin").Error())
	}

	fmt.Printf("configs=%v regex=%v keys=%v", configs, fmt.Sprintf("%T", configs["regex-urls"]), configs["keys"])
	urls, ok := configs["regex-urls"].([]interface{})
	keys, okKeys := configs["keys"].([]interface{})
	if !ok {
		panic(errors.New("incorrect urls").Error())
	}
	if !okKeys {
		panic(errors.New("incorrect keys config").Error())
	}
	urlsRegex := make([]*regexp.Regexp, len(urls))
	for i, u := range urls {
		r, err := regexp.Compile(u.(string))
		if err != nil {
			panic(errors.New(fmt.Sprintf("incorrect url regex err=%v", err)).Error())
		}
		urlsRegex[i] = r
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentPath := r.RequestURI
		mustValidate := false

		for _, regex := range urlsRegex {
			if regex.MatchString(currentPath) {
				mustValidate = true
				break
			}
		}
		if !mustValidate {
			handler.ServeHTTP(w, r)
			return
		}

		keyInformed := r.Header.Get(apiKeyHeader)

		if keyInformed == "" {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		authorized := false
		for _, k := range keys {
			if k.(string) == keyInformed {
				authorized = true
				break
			}
		}

		if !authorized {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)
	}), nil
}
