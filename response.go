package main

type APIResponse struct {
	Data       interface{} `json:"data,omitempty"`
	statusCode int         `json:"-"`
	Message    string      `json:"message"`
}

func NewSuccessResponse(data interface{}, code int, message string) APIResponse {
	return APIResponse{
		Data:       data,
		statusCode: code,
		Message:    message,
	}
}

func NewErrorResponse(code int, message string) APIResponse {
	return APIResponse{
		statusCode: code,
		Message:    message,
	}
}
