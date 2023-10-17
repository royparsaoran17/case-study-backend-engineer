package providererrors

import (
	"fmt"
	"io"
	"net/http"
)

var (
	ErrUserNotFound = NewError("user_not_found", http.StatusBadRequest)
)

type Error struct {
	Code    int    `json:"code"`
	Errors  any    `json:"errors"`
	Message string `json:"message"`
}

func NewError(msg string, code int) Error {
	return Error{
		Code:    code,
		Message: msg,
	}
}

func (e Error) Error() string {
	return e.Message
}

type ErrRequestWithResponse struct {
	Request  *http.Request
	Response *http.Response
}

func NewErrRequestWithResponse(request *http.Request, response *http.Response) ErrRequestWithResponse {
	return ErrRequestWithResponse{
		Request:  request,
		Response: response,
	}
}

func (e ErrRequestWithResponse) Error() string {
	str := fmt.Sprintf("[%d] error request with response from %s %s", e.Response.StatusCode, e.Request.Method, e.Request.URL)
	b, err := io.ReadAll(e.Response.Body)
	if err == nil {
		str = fmt.Sprintf("%s, response: %s", str, string(b))
	}

	return str
}
