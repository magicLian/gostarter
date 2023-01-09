package models

import (
	"errors"
	"fmt"
)

var (
	ErrPermissionDenied = errors.New("permission denied")
	ErrUnauthorized     = errors.New("unauthorized")
)

type HttpErrorResponse struct {
	Message    string
	StatusCode int
	Err        string
}

func (e *HttpErrorResponse) Error() string {
	return e.Err
}

func NewHttpErrorResponse(statusCode int, message string, err error) *HttpErrorResponse {
	if err == nil {
		err = fmt.Errorf(message)
	}
	return &HttpErrorResponse{
		StatusCode: statusCode,
		Message:    message,
		Err:        err.Error(),
	}
}
