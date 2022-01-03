package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type (
	httpHandlerMock struct {
		mockServeHTTP func(http.ResponseWriter, *http.Request)
	}

	writerMock struct {
		registeredHeader  int
		registeredHeaders map[string][]string
	}
)

func (t *httpHandlerMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.mockServeHTTP(w, r)
}

func (w *writerMock) Header() http.Header {
	if w.registeredHeaders == nil {
		w.registeredHeaders = map[string][]string{}
	}
	return w.registeredHeaders
}

func (w *writerMock) Write([]byte) (int, error) {
	return 0, nil
}

func (w *writerMock) WriteHeader(statusCode int) {
	w.registeredHeader = statusCode
}

func TestRegistrable_RegisterHandlers(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		description    string
		extra          map[string]interface{}
		handler        http.Handler
		writer         *writerMock
		request        *http.Request
		validateWriter func(*writerMock) bool
		expectedReturn string
		expectedError  error
	}{
		{
			description: "url not restricted",
			extra: map[string]interface{}{
				"api-key-auth": map[string]interface{}{
					"regex-urls": []interface{}{"^\\/sales-guide\\/api\\/v1\\/companies$"},
					"keys":       []interface{}{"1020"},
				},
			},
			handler: &httpHandlerMock{
				mockServeHTTP: func(writer http.ResponseWriter, request *http.Request) {},
			},
			request: &http.Request{
				RequestURI: "/sales-guide/api/v1/companies/793ad876-5488-437c-84e3-271cc727c9c8",
			},
			writer: &writerMock{},
			validateWriter: func(mock *writerMock) bool {
				return mock.registeredHeader == 0
			},
			expectedError: nil,
		},
		{
			description: "url match and key empty",
			extra: map[string]interface{}{
				"api-key-auth": map[string]interface{}{
					"regex-urls": []interface{}{"^\\/sales-guide\\/api\\/v1\\/companies\\/[0-9a-fA-F]{8}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{12}$"},
					"keys":       []interface{}{"1020"},
				},
			},
			handler: &httpHandlerMock{
				mockServeHTTP: func(writer http.ResponseWriter, request *http.Request) {},
			},
			request: &http.Request{
				RequestURI: "/sales-guide/api/v1/companies/793ad876-5488-437c-84e3-271cc727c9c8",
			},
			writer: &writerMock{},
			validateWriter: func(mock *writerMock) bool {
				return mock.registeredHeader == 401
			},
			expectedError: nil,
		},
		{
			description: "url match and key not match",
			extra: map[string]interface{}{
				"api-key-auth": map[string]interface{}{
					"regex-urls": []interface{}{"^\\/sales-guide\\/api\\/v1\\/companies\\/[0-9a-fA-F]{8}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{12}$"},
					"keys":       []interface{}{"1020", "3040"},
				},
			},
			handler: &httpHandlerMock{
				mockServeHTTP: func(writer http.ResponseWriter, request *http.Request) {},
			},
			request: &http.Request{
				Header: map[string][]string{
					"X-Api-Key": []string{"3041"},
				},
				RequestURI: "/sales-guide/api/v1/companies/793ad876-5488-437c-84e3-271cc727c9c8",
			},
			writer: &writerMock{},
			validateWriter: func(mock *writerMock) bool {
				return mock.registeredHeader == 401
			},
			expectedError: nil,
		},
		{
			description: "url match and key match",
			extra: map[string]interface{}{
				"api-key-auth": map[string]interface{}{
					"regex-urls": []interface{}{"^\\/sales-guide\\/api\\/v1\\/companies\\/[0-9a-fA-F]{8}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{12}$"},
					"keys":       []interface{}{"1020", "3040"},
				},
			},
			handler: &httpHandlerMock{
				mockServeHTTP: func(writer http.ResponseWriter, request *http.Request) {},
			},
			request: &http.Request{
				RequestURI: "/sales-guide/api/v1/companies/793ad876-5488-437c-84e3-271cc727c9c8",
				Header: map[string][]string{
					"X-Api-Key": []string{"3040"},
				},
			},
			writer: &writerMock{},
			validateWriter: func(mock *writerMock) bool {
				return mock.registeredHeader == 0
			},
			expectedError: nil,
		},
		{
			description: "create first user url",
			extra: map[string]interface{}{
				"api-key-auth": map[string]interface{}{
					"regex-urls": []interface{}{"^\\/sales-guide\\/api\\/v1\\/users\\/first-user$"},
					"keys":       []interface{}{"1020", "3040"},
				},
			},
			handler: &httpHandlerMock{
				mockServeHTTP: func(writer http.ResponseWriter, request *http.Request) {},
			},
			request: &http.Request{
				RequestURI: "/sales-guide/api/v1/users/first-user",
				Header: map[string][]string{
					"X-Api-Key": []string{"3040"},
				},
			},
			writer: &writerMock{},
			validateWriter: func(mock *writerMock) bool {
				return mock.registeredHeader == 0
			},
			expectedError: nil,
		},
	}

	for _, c := range cases {
		var r registrable = "api-key-auth"

		h, e := r.registerHandlers(ctx, c.extra, c.handler)

		h.ServeHTTP(c.writer, c.request)

		assert.True(t, c.validateWriter(c.writer))
		assert.Equal(t, c.expectedError, e)
	}
}
