package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// ClientRegisterer is the symbol the plugin loader will try to load. It must implement the RegisterClient interface
var ClientRegisterer = registerer("apikey-plugin")

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
	skeys, ok := extra["allowed-keys"].(string)
	if !ok {
		return nil, errors.New("wrong config, you must inform allowed-keys param")
	}
	keys := strings.Split(skeys, ",")

	// return the actual handler wrapping or your custom logic so it can be used as a replacement for the default http client
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		apikey := req.Header.Get("X-Api-Key")
		validKey := false
		for _, key := range keys {
			if apikey == key {
				validKey = true
				break
			}
		}
		fmt.Println(fmt.Sprintf("****** RESULT OF APIKEY keySent=%v keysRequired=%v approved=%v allHeader=%v", apikey, keys, validKey, req.Header))
		if !validKey {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

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
	fmt.Println("api-key client plugin loaded!!!")
}

func main() {}
