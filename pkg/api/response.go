package api

import (
	"github.com/gin-gonic/gin"
	"time"
)

// APIResponse represents standard API response format
// swagger:model APIResponse
type APIResponse struct {
	Success   bool   `json:"success"`
	Code      int    `json:"code"`
	Timestamp string `json:"timestamp"`
	Data      any    `json:"data,omitempty"`
	Error     string `json:"error,omitempty"`
}

func RespondSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, APIResponse{
		Success:   true,
		Code:      code,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
	})
}

func RespondError(c *gin.Context, code int, msg string) {
	c.JSON(code, APIResponse{
		Success:   false,
		Code:      code,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Error:     msg,
	})
}
