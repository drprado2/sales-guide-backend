package utils

import (
	"encoding/json"
	"github.com/drprado2/react-redux-typescript/internal/domain"
	"github.com/drprado2/react-redux-typescript/pkg/logs"
	"net/http"
)

func HandleError(err error, writter http.ResponseWriter, req *http.Request) {
	if _, ok := err.(*domain.InternalError); ok {
		logs.Logger(req.Context()).WithError(err).Error("create company http request fail with")
		writter.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{
			"error": "ocorreu um erro inesperado, por favor tente novamente.",
		}
		json.NewEncoder(writter).Encode(response)
		return
	}

	response := map[string]string{
		"error": err.Error(),
	}
	writter.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(writter).Encode(response)
	return
}
