package requesthandlers

import (
	"encoding/json"
	"fmt"
	"github.com/drprado2/react-redux-typescript/internal/apptracer"
	"github.com/drprado2/react-redux-typescript/internal/logs"
	"github.com/drprado2/react-redux-typescript/internal/models"
	playerModels "github.com/drprado2/react-redux-typescript/internal/models/players"
	"github.com/drprado2/react-redux-typescript/internal/services/players"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

type UserHandlerInterface interface {
	Create(writter http.ResponseWriter, req *http.Request)
	Update(writter http.ResponseWriter, req *http.Request)
	Delete(writter http.ResponseWriter, req *http.Request)
	GetById(writter http.ResponseWriter, req *http.Request)
	Get(writter http.ResponseWriter, req *http.Request)
}

type UserHandler struct {
	UserService players.UserServiceInterface
	Tracer apptracer.TracerService
}

func (h *UserHandler) Create(writter http.ResponseWriter, req *http.Request) {
	reqBody, _ := ioutil.ReadAll(req.Body)
	player := &playerModels.CreatePlayerRequest{}
	json.Unmarshal(reqBody, player)

	res, err := h.UserService.Create(req.Context(), player)

	writter.Header().Set("Content-Type", "application/json")
	if err != nil {
		writter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writter).Encode(fmt.Sprintf("{error: %v}", err))
		return
	}

	writter.WriteHeader(http.StatusCreated)
	json.NewEncoder(writter).Encode(res)
}

func (h *UserHandler) Update(writter http.ResponseWriter, req *http.Request) {
	playerId := mux.Vars(req)["id"]
	reqBody, _ := ioutil.ReadAll(req.Body)
	player := &playerModels.UpdatePlayerRequest{}
	json.Unmarshal(reqBody, player)

	res, err := h.UserService.Update(req.Context(), playerId, player)

	writter.Header().Set("Content-Type", "application/json")
	if err != nil {
		writter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writter).Encode(fmt.Sprintf("{error: %v}", err))
		return
	}

	writter.WriteHeader(http.StatusOK)
	json.NewEncoder(writter).Encode(res)
}

func (h *UserHandler) Delete(writter http.ResponseWriter, req *http.Request) {
	playerId := mux.Vars(req)["id"]

	err := h.UserService.Delete(req.Context(), playerId)

	writter.Header().Set("Content-Type", "application/json")
	if err != nil {
		writter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writter).Encode(fmt.Sprintf("{error: %v}", err))
		return
	}

	writter.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetById(writter http.ResponseWriter, req *http.Request) {
	playerId := mux.Vars(req)["id"]

	res, err := h.UserService.GetById(req.Context(), playerId)

	writter.Header().Set("Content-Type", "application/json")
	if err != nil {
		writter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writter).Encode(fmt.Sprintf("{error: %v}", err))
		return
	}

	if res == nil {
		writter.WriteHeader(http.StatusNotFound)
		return
	}

	writter.WriteHeader(http.StatusOK)
	json.NewEncoder(writter).Encode(res)
}

func (h *UserHandler) Get(writter http.ResponseWriter, req *http.Request) {
	span, ctx := h.Tracer.SpanFromContext(req.Context())
	defer span.Finish()

	currentPage, err := strconv.Atoi(req.URL.Query().Get("current_page"))
	hasInvalidParameter := false
	if err != nil {
		logs.Logger(ctx).WithError(err).Warnf("Error casting current_page string parameter to int %v", err)
		hasInvalidParameter = true
	}
	itemsByPage, err := strconv.Atoi(req.URL.Query().Get("items_by_page"))
	if err != nil {
		logs.Logger(ctx).WithError(err).Warnf("Error casting items_by_page string parameter to int %v", err)
		hasInvalidParameter = true
	}
	playerId := req.URL.Query().Get("id")
	name := req.URL.Query().Get("name")

	if hasInvalidParameter {
		writter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writter).Encode(models.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	pagParams := &models.PaginationParameters{
		CurrentPage: currentPage,
		ItemsByPage: itemsByPage,
	}
	filter := &playerModels.PlayerFilter{
		ID:   playerId,
		Name: name,
	}

	res, err := h.UserService.GetPaged(ctx, pagParams, filter)

	writter.Header().Set("Content-Type", "application/json")
	if err != nil {
		writter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writter).Encode(models.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if res == nil {
		writter.WriteHeader(http.StatusNotFound)
		return
	}

	writter.WriteHeader(http.StatusOK)
	json.NewEncoder(writter).Encode(res)
}
