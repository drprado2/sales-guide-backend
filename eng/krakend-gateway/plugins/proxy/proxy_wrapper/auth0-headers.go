package main

import (
	"fmt"
	"gopkg.in/square/go-jose.v2/jwt"
	"net/http"
)

const (
	Auth0PluginName     = "auth0-headers"
	AuthHeader          = "Authorization"
	BearerStartPosition = 6
	CompanyIdClaim      = "http://salesguide.com/company-id"
	UserIdClaim         = "http://salesguide.com/user-id"
	EmailClaim          = "http://salesguide.com/email"
	UserIdStartPosition = 6
	UserIdHeader        = "X-User-Id"
	CompanyIdHeader     = "X-Company-Id"
	EmailHeader         = "X-Email"
)

type (
	Auth0Headers struct{}
)

func (a *Auth0Headers) Register(pluginName string, params map[string]interface{}) (func(w http.ResponseWriter, req *http.Request), error) {
	if pluginName != Auth0PluginName {
		return nil, nil
	}

	return a.handleRequest, nil
}

func (a *Auth0Headers) handleRequest(w http.ResponseWriter, req *http.Request) {
	header := req.Header.Get(AuthHeader)
	if header == "" || len(header) < BearerStartPosition {
		return
	}
	var claims map[string]interface{}

	token, err := jwt.ParseSigned(header[BearerStartPosition:])
	if err != nil {
		fmt.Printf("\nerror parsing token, err=%v", err)
		return
	}
	_ = token.UnsafeClaimsWithoutVerification(&claims)

	req.Header.Set(CompanyIdHeader, claims[CompanyIdClaim].(string))
	req.Header.Set(UserIdHeader, claims[UserIdClaim].(string)[UserIdStartPosition:])
	req.Header.Set(EmailHeader, claims[EmailClaim].(string))
	req.Header.Del(AuthHeader)
}
