package routes_api

import (
	"github.com/gin-gonic/gin"
	"github.com/reschedulize/school_course_data"
	"github.com/reschedulize/web_server/db"
	"github.com/reschedulize/web_server/helpers"
	"github.com/reschedulize/web_server/models"
)

type getPlanResponse struct {
	Details       *models.Plan                         `json:"details"`
	CRNToClass    map[string]*school_course_data.Class `json:"crn_to_class"`
	CourseToTypes map[string][]string                  `json:"course_to_types"`
	Schedules     [][]string                           `json:"schedules"`
}

func GETPlanList(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var planIDs []int
	err := db.MySQL.Select(&planIDs, "SELECT `id` FROM `plans` WHERE `user_id` = ?", userID)

	if err != nil {
		helpers.SendInternalServerError(c)
		return
	}

	plans := make([]*models.Plan, len(planIDs))

	for i, id := range planIDs {
		plan, err := models.PlanByID(id)

		if err != nil {
			helpers.SendInternalServerError(c)
			return
		}

		plans[i] = plan
	}

	helpers.SendSuccessfulResponse(c, plans)
}

func GETPlan(c *gin.Context) {
	// Check if plan exists
	planID := c.Param("id")
	userID := c.GetInt64("user_id")

	var id int
	err := db.MySQL.Get(&id, "SELECT `id` FROM `plans` WHERE `id` = ? AND `user_id` = ?", planID, userID)

	if err != nil {
		helpers.SendError(c, 404, "Plan not found")
		return
	}

	// Retrieve plan
	plan, err := models.PlanByID(id)

	if err != nil {
		helpers.SendInternalServerError(c)
		return
	}

	// Retrieve class info for each course
	classes, err := plan.Classes()

	if err != nil {
		helpers.SendInternalServerError(c)
		return
	}

	crnToClass := make(map[string]*school_course_data.Class)
	courseToTypes := make(map[string][]string)

	for _, class := range classes {
		crnToClass[class.CRN] = class

		if _, ok := courseToTypes[class.CourseName]; !ok {
			courseToTypes[class.CourseName] = []string{}
		}

		found := false
		for _, t := range courseToTypes[class.CourseName] {
			if t == class.Type {
				found = true
				break
			}
		}

		if !found {
			courseToTypes[class.CourseName] = append(courseToTypes[class.CourseName], class.Type)
		}
	}

	// Retrieve schedules
	schedules, err := plan.Schedules()

	if err != nil {
		helpers.SendInternalServerError(c)
		return
	}

	helpers.SendSuccessfulResponse(c, getPlanResponse{
		Details:       plan,
		CRNToClass:    crnToClass,
		CourseToTypes: courseToTypes,
		Schedules:     schedules,
	})
}
