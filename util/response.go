package util

import (
	"github.com/gin-gonic/gin"
)

type JSONResponseSuccess struct {
	Data any `json:"data"`
}

type JSONResponseError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func ResponseSuccesJSON(c *gin.Context, data any, statusCode int) {
	c.JSON(statusCode, JSONResponseSuccess{
		Data: data,
	})
}

func ResponseErrorJSON(c *gin.Context, message string, code string, statusCode int) {
	c.JSON(statusCode, JSONResponseError{
		Message: message,
		Code:    code,
	})
}
