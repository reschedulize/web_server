package routes_frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GETLandingPage(c *gin.Context) {
	c.HTML(http.StatusOK, "public.html", gin.H{})
}
