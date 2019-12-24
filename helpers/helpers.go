package helpers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type responseContainer struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
	Error   interface{} `json:"error"`
}

func SendSuccessfulResponse(c *gin.Context, result interface{}) {
	c.JSON(http.StatusOK, responseContainer{
		Success: true,
		Result:  result,
	})
}

func SendError(c *gin.Context, code int, message string) {
	c.JSON(code, responseContainer{
		Success: false,
		Error: gin.H{
			"status_code": code,
			"message":     message,
		},
	})
}

func SendInternalServerError(c *gin.Context) {
	SendError(c, http.StatusInternalServerError, "Internal server error")
}

func SendBadRequestError(c *gin.Context) {
	SendError(c, http.StatusBadRequest, "Bad request")
}