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
	UserController struct {
		serviceManager *domain.ServiceManager
	}
)

func NewUserController(serviceManager *domain.ServiceManager) *UserController {
	return &UserController{serviceManager: serviceManager}
}

func (u *UserController) RegisterRouteHandlers(router *mux.Router) {
	router.
		Path("/v1/users/first-user").
		HandlerFunc(u.CreateFirstUserAction).
		Name("Create First User").
		Methods(Post)
}

func (u *UserController) CreateFirstUserAction(writter http.ResponseWriter, req *http.Request) {
	writter.Header().Set("Content-Type", "application/json")

	reqBody, _ := ioutil.ReadAll(req.Body)
	model := new(usecases.CreateFirstUserInput)
	if err := json.Unmarshal(reqBody, model); err != nil {
		response := map[string]string{
			"error": "paramêtros de entrada inválidos.",
		}
		writter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writter).Encode(response)
		return
	}

	res, err := usecases.CreateFirstUser(req.Context(), u.serviceManager, model)

	if err != nil {
		HandleError(err, writter, req)
		return
	}

	writter.WriteHeader(http.StatusCreated)
	json.NewEncoder(writter).Encode(res)
}
