package utils

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error       bool   `json:"error"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SendError(c *gin.Context, statusCode int, message string) {
	LogError("API_ERROR", map[string]interface{}{
		"status":  statusCode,
		"message": message,
		"path":    c.FullPath(),
	})
	c.JSON(statusCode, ErrorResponse{
		Error:       true,
		ErrorCode:   statusCode,
		Description: message,
	})
}

func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}

func LogError(label string, data map[string]interface{}) {
	output := map[string]interface{}{
		"type":  "error",
		"label": label,
		"data":  data,
	}
	bytes, _ := json.Marshal(output)
	log.Println(string(bytes))
}
