package http

import (
	"encoding/json"
	"github.com/drprado2/sales-guide/internal/domain/errors"
	"net/http"
)

func HandleError(err error, writter http.ResponseWriter, _ *http.Request) {
	if _, ok := err.(*errors.InternalError); ok {
		writter.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{
			"error": "ocorreu um erro inesperado, por favor tente novamente.",
		}
		if err := json.NewEncoder(writter).Encode(response); err != nil {
			panic(err)
		}
		return
	}
	if _, ok := err.(*errors.ConstraintError); ok {
		writter.WriteHeader(http.StatusConflict)
		response := map[string]string{
			"error": err.Error(),
		}
		if err := json.NewEncoder(writter).Encode(response); err != nil {
			panic(err)
		}
		return
	}

	response := map[string]string{
		"error": err.Error(),
	}
	writter.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(writter).Encode(response); err != nil {
		panic(err)
	}
	return
}
