package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Detail)
}

func NewBadRequest(detail string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: "Bad Request", Detail: detail}
}

func NewUnauthorized(detail string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: "Unauthorized", Detail: detail}
}

func NewForbidden(detail string) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: "Forbidden", Detail: detail}
}

func NewNotFound(detail string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: "Not Found", Detail: detail}
}

func NewConflict(detail string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: "Conflict", Detail: detail}
}

func NewInternal(detail string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: "Internal Server Error", Detail: detail}
}

func NewTooManyRequests(detail string) *AppError {
	return &AppError{Code: http.StatusTooManyRequests, Message: "Too Many Requests", Detail: detail}
}
