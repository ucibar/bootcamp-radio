package response

import (
	"encoding/json"
	"net/http"
)

func InvalidCredentials(msg string) *Response {
	return Error(http.StatusUnauthorized, &ErrorPayload{Reason: "invalid_credentials", Message: msg})
}

func InvalidBody(msg string) *Response {
	return Error(http.StatusUnprocessableEntity, &ErrorPayload{Reason: "invalid_body", Message: msg})
}

func InvalidParam(msg string) *Response {
	return Error(http.StatusBadRequest, &ErrorPayload{Reason: "invalid_param", Message: msg})
}

func BadRequest(msg string) *Response {
	return Error(http.StatusBadRequest, &ErrorPayload{Reason: "bad_request", Message: msg})
}

func NotFound(msg string) *Response {
	return Error(http.StatusNotFound, &ErrorPayload{Reason: "not_found", Message: msg})
}

func ServerError() *Response {
	return Error(http.StatusInternalServerError, &ErrorPayload{Reason: "server_error", Message: "something went wrong"})
}

func Success(code int, data interface{}) *Response {
	return &Response{Code: code, Status: "success", Data: data}
}

func Error(code int, error *ErrorPayload) *Response {
	return &Response{Code: code, Status: "error", Error: error}
}

type ErrorPayload struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type Response struct {
	Code   int
	Status string
	Data   interface{}
	Error  *ErrorPayload
}

func NewResponse(code int, status string, data interface{}) *Response {
	return &Response{Code: code, Status: status, Data: data}
}

func (res *Response) MarshalJSON() ([]byte, error) {
	if res.Error != nil {
		return json.Marshal(&struct {
			Status string        `json:"status"`
			Error  *ErrorPayload `json:"error"`
		}{
			Status: res.Status,
			Error:  res.Error,
		})
	} else {
		return json.Marshal(&struct {
			Status string      `json:"status"`
			Data   interface{} `json:"data, omitempty"`
		}{
			Status: res.Status,
			Data:   res.Data,
		})
	}
}
