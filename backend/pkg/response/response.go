package response

import (
	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    any `json:"data,omitempty"`
	Error   any `json:"error,omitempty"`
}

func Success(c *gin.Context, code int, data any) {
	c.JSON(code, ApiResponse{
		Status: "success",
		Data:   data,
	})
}

func Error(c *gin.Context, code int, message string, err error) {
    var errMsg any
    if err != nil {
        errMsg = err.Error()
    } else {
        errMsg = nil
    }
	c.JSON(code, ApiResponse{
		Status:  "error",
		Message: message,
		Error:   errMsg,
	})
}
