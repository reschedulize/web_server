package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/reschedulize/web_server/db"
)

func isCookieValid(c *gin.Context) (bool, int64, string) {
	sessionID, err := c.Cookie("SID")

	if err != nil {
		return false, 0, ""
	}

	var userID int64
	err = db.MySQL.Get(&userID, "SELECT `user_id` FROM `tokens` WHERE `session_id` = ?", sessionID)

	if err != nil {
		return false, 0, ""
	}

	return true, userID, sessionID
}
