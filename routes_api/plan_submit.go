package routes_api

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/reschedulize/web_server/db"
	"github.com/reschedulize/web_server/helpers"
	"github.com/reschedulize/web_server/models"
	"sort"
)

var invalidInputError = errors.New("invalid input error")

func verifyInput(plan *models.Plan) error {
	// Check term is valid
	if len(plan.Term) != 6 {
		return invalidInputError
	}

	courses, err := db.UCRAPI.Courses(plan.Term, 10000)

	if err != nil {
		return err
	}

	if len(courses) == 0 {
		return invalidInputError
	}

	// TODO: Set limits on number of courses

	// Check courses are valid
	seen := make(map[string]bool)

	for _, group := range plan.CourseGroups {
		for _, course := range group {
			_, ok := seen[course]

			// Check for duplicate courses
			if ok {
				return invalidInputError
			} else {
				seen[course] = true
			}

			i := sort.SearchStrings(courses, course)

			if i >= len(courses) || courses[i] != course {
				return invalidInputError
			}
		}
	}

	return nil
}

func POSTPlan(c *gin.Context) {
	var plan models.Plan
	err := c.BindJSON(&plan)

	if err != nil {
		helpers.SendBadRequestError(c)
		return
	}

	err = verifyInput(&plan)

	if err != nil {
		if err == invalidInputError {
			helpers.SendBadRequestError(c)
			return
		} else {
			helpers.SendInternalServerError(c)
			return
		}
	}

	userID := c.GetInt64("user_id")

	//ok, err = consumeCredit(userID)
	//
	//if err != nil {
	//	helpers.SendInternalServerError(c)
	//	return
	//}
	//
	//if !ok {
	//	helpers.SendError(c, http.StatusTooManyRequests, "No credits remaining")
	//	return
	//}

	// Any error beyond this point will refund a credit to the user
	//defer func() {
	//	_, _ = db.MySQL.Exec("UPDATE `users` SET `remaining_plans` = `remaining_plans` + 1 WHERE `id` = ?", userID)
	//}()

	planID, err := models.SavePlan(&plan, userID)

	if err != nil {
		helpers.SendInternalServerError(c)
		return
	}

	helpers.SendSuccessfulResponse(c, gin.H{
		"id": planID,
	})
}
