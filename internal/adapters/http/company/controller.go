package company

import (
	"encoding/json"
	"github.com/drprado2/react-redux-typescript/internal/adapters/http/utils"
	"github.com/drprado2/react-redux-typescript/internal/domain/usecases"
	"github.com/drprado2/react-redux-typescript/pkg/httpserver"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func RegisterRouteHandlers(router *mux.Router) {
	router.
		Path("/v1/companies").
		HandlerFunc(CreateCompanyAction).
		Name("Create Company").
		Methods(httpserver.Post)

	router.
		Path("/v1/companies/{id}").
		HandlerFunc(GetCompanyByIdAction).
		Name("Get Company").
		Methods(httpserver.Get)
}

func CreateCompanyAction(writter http.ResponseWriter, req *http.Request) {
	writter.Header().Set("Content-Type", "application/json")

	reqBody, _ := ioutil.ReadAll(req.Body)
	model := new(usecases.CreateCompanyInput)
	if err := json.Unmarshal(reqBody, model); err != nil {
		response := map[string]string{
			"error": err.Error(),
		}
		writter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writter).Encode(response)
		return
	}

	res, err := usecases.CreateCompany(req.Context(), model)

	if err != nil {
		utils.HandleError(err, writter, req)
		return
	}

	writter.WriteHeader(http.StatusCreated)
	json.NewEncoder(writter).Encode(res)
}

func GetCompanyByIdAction(writter http.ResponseWriter, req *http.Request) {
	writter.Header().Set("Content-Type", "application/json")

	params := mux.Vars(req)
	res, err := usecases.GetCompanyByID(req.Context(), params["id"])

	if err != nil {
		utils.HandleError(err, writter, req)
		return
	}

	writter.WriteHeader(http.StatusCreated)
	json.NewEncoder(writter).Encode(res)
}
