package routes_api

import (
	"github.com/gin-gonic/gin"
	"github.com/reschedulize/school_course_data"
	"github.com/reschedulize/web_server/db"
	"github.com/reschedulize/web_server/helpers"
	"strings"
)

type getClassesResponse struct {
	Classes []*school_course_data.Class
	Error   error
}

func GETTerms(c *gin.Context) {
	terms, err := db.UCRAPI.Terms(4)

	if err != nil {
		helpers.SendInternalServerError(c)
		return
	}

	helpers.SendSuccessfulResponse(c, terms)
}

func GETCourseList(c *gin.Context) {
	term := c.Param("term")

	if len(term) != 6 {
		helpers.SendBadRequestError(c)
		return
	}

	courses, err := db.UCRAPI.Courses(term, 10000)

	if err != nil {
		helpers.SendInternalServerError(c)
		return
	}

	if len(courses) == 0 {
		helpers.SendBadRequestError(c)
		return
	}

	helpers.SendSuccessfulResponse(c, courses)
}

func GETCourseStats(c *gin.Context) {
	term := c.Param("term")
	course := c.Param("course")
	classes, err := db.UCRAPI.Classes(term, course, 100)

	if err != nil {
		helpers.SendInternalServerError(c)
		return
	}

	var units int
	var minutes uint16
	var types []string

	seenTypes := make(map[string]bool)

	for _, class := range classes {
		_, ok := seenTypes[class.Type]

		if !ok {
			seenTypes[class.Type] = true

			units += class.Units
			types = append(types, class.Type)

			//for _, day := range class.Schedule {
			//	for _, timeRange := range day {
			//		minutes += timeRange.End - timeRange.Begin
			//	}
			//}
		}
	}

	helpers.SendSuccessfulResponse(c, gin.H{
		"units":   units,
		"minutes": minutes,
		"types":   types,
	})
}

func GETClasses(c *gin.Context) {
	term := c.Param("term")
	courses := strings.Split(c.Param("courses"), ",")

	var result []*school_course_data.Class

	ch := make(chan getClassesResponse, len(courses))

	for _, course := range courses {
		go func(course string) {
			classes, err := db.UCRAPI.Classes(term, course, 100)

			ch <- getClassesResponse{
				Classes: classes,
				Error:   err,
			}
		}(course)
	}

	for i := 0; i < len(courses); i++ {
		response := <-ch

		if response.Error != nil {
			helpers.SendInternalServerError(c)
			return
		}

		result = append(result, response.Classes...)
	}

	helpers.SendSuccessfulResponse(c, result)
}
