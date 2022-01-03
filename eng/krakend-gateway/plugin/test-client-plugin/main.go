package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// ClientRegisterer is the symbol the plugin loader will try to load. It must implement the RegisterClient interface
var ClientRegisterer = registerer("krakend-example")

var client = &http.Client{Timeout: 3 * time.Second}

type registerer string

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(string(r), r.registerClients)
}

func (r registerer) registerClients(ctx context.Context, extra map[string]interface{}) (http.Handler, error) {
	// check the passed configuration and initialize the plugin
	name, ok := extra["name"].(string)
	if !ok {
		return nil, errors.New("wrong config")
	}
	if name != string(r) {
		return nil, fmt.Errorf("unknown register %s", name)
	}

	// return the actual handler wrapping or your custom logic so it can be used as a replacement for the default http client
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		req.Header.Set("x-perereca", "zuado123")
		rs, err := client.Do(req)
		if err != nil {
			fmt.Println(fmt.Sprintf("****** ERROR %v", err))
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}
		defer rs.Body.Close()

		rsBodyBytes, err := ioutil.ReadAll(rs.Body)
		if err != nil {
			fmt.Println(fmt.Sprintf("****** ERROR %v", err))
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}

		fmt.Println(fmt.Sprintf("****** RUNNING WITH SUCCESS"))

		for name, values := range rs.Header {
			for _, value := range values {
				w.Header().Set(name, value)
			}
		}
		w.Header().Set("x-toba", "123")
		w.WriteHeader(rs.StatusCode)
		w.Write(rsBodyBytes)
	}), nil
}

func init() {
	fmt.Println("krakend-example client plugin loaded!!!")
}

func main() {}
