package main

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
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

	mockHttpClient struct {
		reqCalled *http.Request
	}

	mockReadCloser struct {
	}
)

func (m *mockReadCloser) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (m *mockReadCloser) Close() error {
	return nil
}

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	m.reqCalled = req
	resp := &http.Response{}
	resp.Body = &mockReadCloser{}
	return resp, nil
}

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

func TestRegisterer_RegisterClients(t *testing.T) {
	reg := registerer("proxy-proxy_wrapper-plugin")
	ctx := context.Background()
	extra := map[string]interface{}{}
	mockClient := &mockHttpClient{}
	client = mockClient

	t.Run("without plugins", func(t *testing.T) {
		_, err := reg.registerClients(ctx, extra)
		assert.Error(t, err)
	})

	t.Run("unavailable plugin", func(t *testing.T) {
		configs := []*ProxyPluginRegister{
			{
				Name:   "no-auth0-headers",
				Params: nil,
			},
		}
		configsJson, err := json.Marshal(configs)
		assert.NoError(t, err)

		extra["plugins"] = string(configsJson)
		_, err = reg.registerClients(ctx, extra)
		assert.NoError(t, err)
		assert.Len(t, handlers, 0)
	})

	t.Run("auth0 plugin", func(t *testing.T) {
		configs := []*ProxyPluginRegister{
			{
				Name:   "auth0-headers",
				Params: nil,
			},
		}
		configsJson, err := json.Marshal(configs)
		assert.NoError(t, err)

		extra["plugins"] = string(configsJson)

		handler, err := reg.registerClients(ctx, extra)
		assert.NoError(t, err)
		request := &http.Request{
			RequestURI: "/sales-guide/api/v1/companies/793ad876-5488-437c-84e3-271cc727c9c8",
			Header: map[string][]string{
				"Authorization": []string{"Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlpGYlByVHdGUW5Oa2k4SkJGTDRVOCJ9.eyJodHRwOi8vc2FsZXNndWlkZS5jb20vY29tcGFueS1pZCI6ImZkOGY5ZDZmLTg1ZDgtNGMxNy1hMzNkLTMwMTcxZDkxZTk5ZCIsImh0dHA6Ly9zYWxlc2d1aWRlLmNvbS91c2VyLWlkIjoiYXV0aDB8NjBmMGViMWY2MTBhNzYwMDY5ZWJkMDg1IiwiaHR0cDovL3NhbGVzZ3VpZGUuY29tL2VtYWlsIjoiZHJwcmFkbzJAZ21haWwuY29tIiwiaXNzIjoiaHR0cHM6Ly9kcnByYWRvMi51cy5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NjBmMGViMWY2MTBhNzYwMDY5ZWJkMDg1IiwiYXVkIjpbImh0dHA6Ly9sb2NhbGhvc3Q6ODAwMCIsImh0dHBzOi8vZHJwcmFkbzIudXMuYXV0aDAuY29tL3VzZXJpbmZvIl0sImlhdCI6MTY0MDg5ODcxNSwiZXhwIjoxNjQwOTg1MTE1LCJhenAiOiJBUUU0VHgxWXlrYm1henlicVJpMHJrZDcxRlRKbzB2WiIsInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwgb2ZmbGluZV9hY2Nlc3MiLCJwZXJtaXNzaW9ucyI6WyJyZWFkOmFsbCJdfQ.JKtyoHIglH3AIiXNCmT7uIZ_Zfyl01Y5AsiwzThDutYSkIkIT5MZRmpCoE-E8VzpYUZZgZFpmTodXB9vRiQStH_y49Q9Bw8kO8qIxiJo-WZYEOCaPd_yubj3bv_Bp67yEvSZn10RteKhNFf4SsxWe2bYIInoWTYA8m8LAsHCCI0U_TSMe_azdaoY_Ciu_Hn9ZuoUcxtMWpaFkzkoF_Io7nJFRRqDVyRX3Oh9v2GmGuyWj0cSwyGjPhRwFc45u8s3xtrHtQ96R-zuqXvC7KYrbxtI9F-bPaa-Q-AbDK_uYYsVGcwdqGuQFrTjFC-7EVnaZ_8_DWDx6I9FQQrUOEZ4jg"},
			},
		}
		writer := &writerMock{}

		handler.ServeHTTP(writer, request)

		assert.Equal(t, mockClient.reqCalled.Header.Get("Authorization"), "")
		assert.Equal(t, mockClient.reqCalled.Header.Get("X-Company-Id"), "fd8f9d6f-85d8-4c17-a33d-30171d91e99d")
		assert.Equal(t, mockClient.reqCalled.Header.Get("X-User-Id"), "60f0eb1f610a760069ebd085")
		assert.Equal(t, mockClient.reqCalled.Header.Get("X-Email"), "drprado2@gmail.com")
	})
}
