package models

type GenericResponse struct {
	Status  int    `json:"status"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func NewGenericResponse(status int, data any, message string) *GenericResponse {
	return &GenericResponse{
		Status:  status,
		Data:    data,
		Message: message,
	}
}
