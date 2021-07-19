package adapters

import (
	"encoding/json"
	"github.com/drprado2/react-redux-typescript/internal/domain"
	"github.com/drprado2/react-redux-typescript/internal/domain/usecases"
	"github.com/drprado2/react-redux-typescript/pkg/httpserver"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type (
	CompanyHttpAdapter struct {
		logger domain.Logger
	}
)

func NewCompanyHttpAdapter(logger domain.Logger) *CompanyHttpAdapter {
	return &CompanyHttpAdapter{
		logger: logger,
	}
}

func (cha *CompanyHttpAdapter) RegisterRouteHandlers(router *mux.Router) {
	router.
		Path("/v1/players").
		HandlerFunc(cha.CreateCompanyAction).
		Name("Create Company").
		Methods(httpserver.Post)
}

func (cha *CompanyHttpAdapter) CreateCompanyAction(writter http.ResponseWriter, req *http.Request) {
	reqBody, _ := ioutil.ReadAll(req.Body)
	model := &usecases.CreateCompanyInput{}
	if err := json.Unmarshal(reqBody, model); err != nil {
		response := map[string]interface{}{
			"error": err,
		}
		resp, _ := json.Marshal(response)
		writter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writter).Encode(resp)
		return
	}

	res, err := usecases.CreateCompany(req.Context(), model)

	writter.Header().Set("Content-Type", "application/json")
	if err != nil {
		if _, ok := err.(*domain.InternalError); ok {
			cha.logger.ErrorWithFieldsf(req.Context(), map[string]interface{}{"error": err}, "create company http request fail with error=%v", err)
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

	writter.WriteHeader(http.StatusCreated)
	json.NewEncoder(writter).Encode(res)
}
