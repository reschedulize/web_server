package routes_frontend

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/reschedulize/web_server/db"
	"github.com/reschedulize/web_server/helpers"
	"net/http"
)

type LoginForm struct {
	Key             string `json:"key" binding:"required"`
	CaptchaResponse string `json:"captcha_response" binding:"required"`
}

func GETLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "public.html", gin.H{})
}

func POSTLogin(c *gin.Context) {
	var form LoginForm
	err := c.BindJSON(&form)

	if err != nil {
		helpers.SendError(c, http.StatusBadRequest, "Bad request")
		return
	}

	// Verify key
	row := db.MySQL.QueryRow("SELECT `id` FROM `users` WHERE `key` = ?", form.Key)

	var id int64
	err = row.Scan(&id)

	if err != nil {
		helpers.SendError(c, http.StatusUnauthorized, "Invalid key")
		return
	}

	// Generate session id
	var sessionID string

	for {
		b := make([]byte, 32)
		_, err = rand.Read(b)

		if err != nil {
			helpers.SendInternalServerError(c)
			return
		}

		sessionID = hex.EncodeToString(b)
		_, err := db.MySQL.Exec("INSERT INTO `tokens` (`user_id`, `session_id`) VALUES (?, ?)", id, sessionID)

		if err != nil {
			me, ok := err.(*mysql.MySQLError)

			// Repeat loop if there is a duplicate session id
			if !ok || me.Number != 1062 {
				helpers.SendInternalServerError(c)
				return
			}
		} else {
			break
		}
	}

	c.SetCookie("SID", sessionID, 2147483647, "/", "144.202.126.188", false, false)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func GETLogout(c *gin.Context) {
	_, err := db.MySQL.Exec("DELETE FROM `tokens` WHERE `session_id` = ?", c.MustGet("session_id"))

	if err != nil {
		helpers.SendInternalServerError(c)
		return
	}

	c.SetCookie("SID", "", 0, "/", "144.202.126.188", false, false)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
