package util

import (
	"github.com/gin-gonic/gin"
)

// Response ...
type Response struct {
	Message    string      `json:"message,omitempty"`
	StatusCode int         `json:"statusCode,omitempty"`
	Code       string      `json:"code,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

// ResponseError ...
type ResponseError struct {
	Message    string        `json:"message,omitempty"`
	StatusCode int           `json:"statusCode,omitempty"`
	Code       string        `json:"code,omitempty"`
	Errors     []interface{} `json:"errors,omitempty"`
}

// Send ResponseError to HTTP in JSON format
func (r ResponseError) Send(c *gin.Context, code int) {
	c.JSON(code, r)
}
