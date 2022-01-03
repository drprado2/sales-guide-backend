package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	registerer  string
	ProxyPlugin interface {
		Register(pluginName string, params map[string]interface{}) (func(w http.ResponseWriter, req *http.Request), error)
	}
	ProxyPluginRegister struct {
		Name   string                 `json:"name,omitempty"`
		Params map[string]interface{} `json:"params,omitempty"`
	}
	HttpClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
)

var (
	// ClientRegisterer is the symbol the plugin loader will try to load. It must implement the RegisterClient interface
	ClientRegisterer = registerer("proxy_wrapper")

	client HttpClient = &http.Client{Timeout: 6 * time.Second}

	availablePlugins = []ProxyPlugin{
		&Auth0Headers{},
	}

	handlers = make([]func(w http.ResponseWriter, req *http.Request), 0)
)

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(string(r), r.registerClients)
}

func (r registerer) registerClients(ctx context.Context, extra map[string]interface{}) (http.Handler, error) {
	plugins, ok := extra["plugins"].(string)
	if !ok {
		return nil, errors.New("[PROXY WRAPPER] wrong plugins config")
	}
	fmt.Printf("[PROXY WRAPPER] configs=%v", plugins)

	configs := make([]*ProxyPluginRegister, 0)
	if err := json.Unmarshal([]byte(plugins), &configs); err != nil {
		panic(err)
	}

	for _, p := range configs {
		fmt.Printf("[PROXY WRAPPER] configs=%v plugin=%v type=%v", plugins, p, fmt.Sprintf("%T", p))
		for _, pp := range availablePlugins {
			handler, err := pp.Register(p.Name, p.Params)
			if err != nil {
				panic(errors.New(fmt.Sprintf("invalid plugin name=%s err=%v", p.Name, err)).Error())
			}
			if handler != nil {
				handlers = append(handlers, handler)
			}
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		for _, h := range handlers {
			h(w, req)
		}

		rs, err := client.Do(req)
		if err != nil {
			fmt.Println(fmt.Sprintf("ERROR %v", err))
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

		for name, values := range rs.Header {
			for _, value := range values {
				w.Header().Set(name, value)
			}
		}
		w.WriteHeader(rs.StatusCode)
		w.Write(rsBodyBytes)
	}), nil
}

func init() {
	fmt.Println("api-key client plugin loaded!!!")
}

func main() {}
