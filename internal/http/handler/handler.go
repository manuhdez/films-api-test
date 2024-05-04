package handler

import (
	"errors"
)

var (
	ErrUnauthorized = errors.New("request is not unauthorized")
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(e error) ErrorResponse {
	return ErrorResponse{Message: e.Error()}
}
