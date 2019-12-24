package routes_api

import (
	"github.com/gin-gonic/gin"
	"github.com/reschedulize/web_server/db"
	"github.com/reschedulize/web_server/helpers"
)

func GETUser(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var remainingPlans uint
	err := db.MySQL.Get(&remainingPlans, "SELECT `remaining_plans` FROM `users` WHERE `id` = ?", userID)

	if err != nil {
		helpers.SendInternalServerError(c)
		return
	}

	helpers.SendSuccessfulResponse(c, gin.H{
		"remaining_plans": remainingPlans,
	})
}
