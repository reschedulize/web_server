package routes_frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GETViewAllPlans(c *gin.Context) {
	c.HTML(http.StatusOK, "private.html", gin.H{})
}

func GETViewPlan(c *gin.Context) {
	c.HTML(http.StatusOK, "private.html", gin.H{})
}

func GETNewPlan(c *gin.Context) {
	c.HTML(http.StatusOK, "private.html", gin.H{})
}
