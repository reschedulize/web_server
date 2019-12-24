package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Public(c *gin.Context) {
	valid, _, _ := isCookieValid(c)

	if !valid {
		c.Next()
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/plans")
}
