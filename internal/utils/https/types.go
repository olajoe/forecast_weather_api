package https

import (
	"net/http"
)

const (
	HeaderContentTypeKey   = "Content-Type"
	HeaderAuthorizationKey = "Authorization"
)

type ProxyRequest struct {
	Method      string
	URL         string
	QueryParams map[string]string
	PathParams  map[string]string
	Headers     map[string]string
	Body        interface{}
}

type ProxyResponse struct {
	StatusCode int
	Data       interface{}
}

type ErrorResponse struct {
	Status  int         `json:"status"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

func (e ErrorResponse) Error() string {
	if e.Message == "" {
		return e.Code
	}

	return e.Message
}

func NewErrorResponse(status int, code string, message string) ErrorResponse {
	return ErrorResponse{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func NewErrorResponseBadRequest(err error) ErrorResponse {
	return NewErrorResponse(http.StatusBadRequest, "bad-request", err.Error())
}

func NewErrorResponseUnauthorized(err error) ErrorResponse {
	return NewErrorResponse(http.StatusUnauthorized, "unauthorized", err.Error())
}

func NewErrorResponseForbidden(err error) ErrorResponse {
	return NewErrorResponse(http.StatusForbidden, "forbidden", err.Error())
}

func NewErrorResponseConflict(err error) ErrorResponse {
	return NewErrorResponse(http.StatusConflict, "conflict", err.Error())
}

func NewErrorResponseUnprocessableEntity(err error) ErrorResponse {
	return NewErrorResponse(http.StatusUnprocessableEntity, "unprocessable-entity", err.Error())
}

func NewErrorResponseInternalServerError(err error) ErrorResponse {
	return NewErrorResponse(http.StatusInternalServerError, "internal-server-error", err.Error())
}
