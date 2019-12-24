package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/reschedulize/web_server/db"
	"net/http"
	"time"
)

func Private(c *gin.Context) {
	valid, userID, sessionID := isCookieValid(c)

	if !valid {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	go db.MySQL.Exec("UPDATE tokens SET `last_updated` = ? WHERE `session_id` = ?", time.Now(), sessionID)

	c.Set("user_id", userID)
	c.Set("session_id", sessionID)
	c.Next()
}
