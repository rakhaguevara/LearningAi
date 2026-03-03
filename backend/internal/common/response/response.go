package response

import (
	"net/http"

	apperr "github.com/adaptive-ai-learn/backend/internal/common/errors"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
}

type ErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Success: true, Data: data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{Success: true, Data: data})
}

func Err(c *gin.Context, err error) {
	if appErr, ok := err.(*apperr.AppError); ok {
		c.JSON(appErr.Code, Response{
			Success: false,
			Error: &ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
				Detail:  appErr.Detail,
			},
		})
		return
	}
	// Expose the actual error message so the frontend can show it clearly.
	// In production, consider filtering sensitive details via an env flag.
	msg := "Internal Server Error"
	if err != nil {
		msg = err.Error()
	}
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error: &ErrorBody{
			Code:    http.StatusInternalServerError,
			Message: msg,
		},
	})
}

func ErrStatus(c *gin.Context, statusCode int, err error) {
	msg := "Internal Server Error"
	if err != nil {
		msg = err.Error()
	}
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorBody{
			Code:    statusCode,
			Message: msg,
		},
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
