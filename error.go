package vodka

import (
	"runtime"
)

// Error - framework error type
type Error struct {
	httpCode int
	Message  string      `json:"message"`
	Stack    interface{} `json:"stack,omitempty"`
	Info     interface{} `json:"info"`
}

// NewServerError - 500 error decorator
func NewServerError(message string, info interface{}) error {
	return NewError(ErrorServerErrorCode, message, info)
}

// NewBadRequestError - 400 error decorator
func NewBadRequestError(message string, info interface{}) error {
	return NewError(ErrorBadRequestCode, message, info)
}

// NewUnathorizedError - 401 error decorator
func NewUnathorizedError(message string, info interface{}) error {
	return NewError(ErrorUnathorizedCode, message, info)
}

// NewError - Error constructor
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
