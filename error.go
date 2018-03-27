package vodka

import (
	"runtime"
)

const (
	// BadRequestCode - server HTTP code for BadRequest 400
	BadRequestCode = 400
	// UnathorizedCode - server HTTP code for Unathorized 401
	UnathorizedCode = 401
	// ServerErrorCode - server HTTP code for ServerError 500
	ServerErrorCode = 500
)

type Error struct {
	httpCode int
	Message  string      `json:"message"`
	Stack    interface{} `json:"stack,omitempty"`
	Info     interface{} `json:"info"`
}

func NewServerError(message string, info interface{}) error {
	return NewError(ServerErrorCode, message, info)
}

func NewBadRequestError(message string, info interface{}) error {
	return NewError(BadRequestCode, message, info)
}

func NewUnathorizedError(message string, info interface{}) error {
	return NewError(UnathorizedCode, message, info)
}

func NewError(httpCode int, message string, info interface{}) error {
	buf := make([]byte, 2048)
	runtime.Stack(buf, true)
	return Error{
		httpCode: httpCode,
		Message:  message,
		Info:     info,
		Stack:    string(buf),
	}
}

func (e Error) Error() string {
	return e.Message
}
