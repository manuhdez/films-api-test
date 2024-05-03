package handler

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(e error) ErrorResponse {
	return ErrorResponse{Message: e.Error()}
}
