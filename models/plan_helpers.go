package models

import (
	"errors"
	"github.com/reschedulize/web_server/db"
	"sort"
)

var UnableToRetrievePlanError = errors.New("unable to retrieve plan")

func PlanByID(id int) (*Plan, error) {
	var plan Plan

	err := db.MySQL.Get(&plan, "SELECT `id`, `term` FROM `plans` WHERE `id` = ?", id)

	if err != nil {
		return nil, UnableToRetrievePlanError
	}

	// Course groups
	var courseGroupIDs []uint
	err = db.MySQL.Select(&courseGroupIDs, "SELECT `id` FROM `course_groups` WHERE `plan_id` = ?", plan.ID)

	if err != nil {
		return nil, UnableToRetrievePlanError
	}

	for _, courseGroupID := range courseGroupIDs {
		var courses []string
		err = db.MySQL.Select(&courses, "SELECT `course_name` FROM `course_group_courses` WHERE `course_group_id` = ?", courseGroupID)

		if err != nil {
			return nil, UnableToRetrievePlanError
		}

		plan.CourseGroups = append(plan.CourseGroups, courses)
	}

	if len(courseGroupIDs) == 0 {
		plan.CourseGroups = [][]string{}
	}

	return &plan, nil
}

func SavePlan(plan *Plan, userID int64) (int64, error) {
	// Sort group
	for i := range plan.CourseGroups {
		sort.Strings(plan.CourseGroups[i])
	}

	// Sort groups
	sort.Slice(plan.CourseGroups, func(i, j int) bool {
		return sort.StringsAreSorted([]string{plan.CourseGroups[i][0], plan.CourseGroups[j][0]})
	})

	// Create plan
	rows, err := db.MySQL.Exec("INSERT INTO `plans` (`user_id`, `term`) VALUES (?, ?)", userID, plan.Term)

	if err != nil {
		return 0, err
	}

	planID, err := rows.LastInsertId()

	if err != nil {
		return 0, err
	}

	// Course groups
	for _, group := range plan.CourseGroups {
		if len(group) == 0 {
			continue
		}

		rows, err := db.MySQL.Exec("INSERT INTO `course_groups` (`plan_id`) VALUES (?)", planID)

		if err != nil {
			return 0, err
		}

		groupID, err := rows.LastInsertId()

		if err != nil {
			return 0, err
		}

		for _, course := range group {
			_, err := db.MySQL.Exec("INSERT INTO `course_group_courses` (`course_group_id`, `course_name`) VALUES (?, ?)", groupID, course)

			if err != nil {
				return 0, err
			}
		}
	}

	return planID, nil
}
