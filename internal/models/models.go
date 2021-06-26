package models

import "net/http"

type PaginationParameters struct {
	CurrentPage int `json:"current_page"`
	ItemsByPage int `json:"items_by_page"`
}

type PaginationResponse struct {
	Data       interface{} `json:"data"`
	TotalItems int         `json:"total_items"`
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (rec *StatusRecorder) WriteHeader(code int) {
	rec.Status = code
	rec.ResponseWriter.WriteHeader(code)
}
