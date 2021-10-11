package user

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
		Path("/v1/users/first-user").
		HandlerFunc(CreateFirstUserAction).
		Name("Create First User").
		Methods(httpserver.Post)
}

func CreateFirstUserAction(writter http.ResponseWriter, req *http.Request) {
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

	res, err := usecases.CreateFirstUser(req.Context(), model)

	if err != nil {
		utils.HandleError(err, writter, req)
		return
	}

	writter.WriteHeader(http.StatusCreated)
	json.NewEncoder(writter).Encode(res)
}
