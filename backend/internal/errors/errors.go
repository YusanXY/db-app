package errors

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewBadRequestError(message string) *AppError {
	return &AppError{Code: 400, Message: message}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{Code: 401, Message: message}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{Code: 403, Message: message}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{Code: 404, Message: message}
}

func NewInternalError(message string) *AppError {
	return &AppError{Code: 500, Message: message}
}

func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		c.JSON(appErr.Code, gin.H{
			"code":    appErr.Code,
			"message": appErr.Message,
		})
	} else {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
	}
}

