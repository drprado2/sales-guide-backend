package http

import (
	"encoding/json"
	"github.com/drprado2/sales-guide/internal/domain"
	"github.com/drprado2/sales-guide/internal/domain/usecases"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type (
	CompanyController struct {
		serviceManager *domain.ServiceManager
	}
)

func NewCompanyController(serviceManager *domain.ServiceManager) *CompanyController {
	return &CompanyController{serviceManager: serviceManager}
}

func (c *CompanyController) RegisterRouteHandlers(router *mux.Router) {
	router.
		Path("/v1/companies").
		HandlerFunc(c.CreateAction).
		Name("Create Company").
		Methods(Post)

	router.
		Path("/v1/companies/{id}").
		HandlerFunc(c.GetByIdAction).
		Name("Get Company").
		Methods(Get)
}

func (c *CompanyController) CreateAction(writter http.ResponseWriter, req *http.Request) {
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

	res, err := usecases.CreateCompany(req.Context(), c.serviceManager, model)

	if err != nil {
		HandleError(err, writter, req)
		return
	}

	writter.WriteHeader(http.StatusCreated)
	json.NewEncoder(writter).Encode(res)
}

func (c *CompanyController) GetByIdAction(writter http.ResponseWriter, req *http.Request) {
	writter.Header().Set("Content-Type", "application/json")

	params := mux.Vars(req)
	res, err := usecases.GetCompanyByID(req.Context(), c.serviceManager, params["id"])

	if err != nil {
		HandleError(err, writter, req)
		return
	}

	writter.WriteHeader(http.StatusOK)
	json.NewEncoder(writter).Encode(res)
}
