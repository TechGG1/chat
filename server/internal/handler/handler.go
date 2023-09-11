package handler

import (
	"encoding/json"
	"errors"
	"github.com/TechGG1/chat/server/internal/logging"
	"github.com/TechGG1/chat/server/internal/service"
	"net/http"
)

type Handler struct {
	service *service.Service
	logger  *logging.Logger
}

func NewHandler(service *service.Service, logger *logging.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

var (
	ErrInvalidCredentials  = errors.New("invalid login credentials")
	ErrInRequestMarshaling = errors.New("invalid/bad request paramenters")
	ErrDuplicateEmail      = errors.New("email already exists")
	ErrMalformedToken      = errors.New("malformed jwt token")
)

type ErrorResponse struct {
	Message string `json:"Message"`
	Code    int    `json:"Code"`
	Status  bool   `json:"Status"`
}

func ErrResponse(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	errCode := codeFrom(err)
	w.WriteHeader(errCode)
	res := ErrorResponse{Message: err.Error(), Status: false, Code: errCode}
	data, err := json.Marshal(res)
	if err != nil {
		return
	}
	w.Write(data)
}

func codeFrom(err error) int {
	switch err {
	case ErrInvalidCredentials:
		return http.StatusBadRequest
	case ErrDuplicateEmail:
		return http.StatusBadRequest
	case ErrInRequestMarshaling:
		return http.StatusBadRequest
	case ErrInRequestMarshaling:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
