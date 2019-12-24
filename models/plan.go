package models

import (
	"github.com/reschedulize/algorithm"
	"github.com/reschedulize/school_course_data"
	"github.com/reschedulize/web_server/db"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Plan struct {
	ID           int64      `db:"id" json:"id"`
	Term         string     `db:"term" json:"term"`
	CourseGroups [][]string `json:"course_groups"`
}

func (p *Plan) Classes() ([]*school_course_data.Class, error) {
	courses := make(map[string]bool)

	for _, group := range p.CourseGroups {
		for _, course := range group {
			courses[course] = true
		}
	}

	var result []*school_course_data.Class
	var err error

	var wg sync.WaitGroup
	var mutex sync.Mutex

	for course := range courses {
		wg.Add(1)

		go func(course string) {
			defer wg.Done()

			var classes []*school_course_data.Class
			classes, err = db.UCRAPI.Classes(p.Term, course, 100)

			if err != nil {
				return
			}

			mutex.Lock()
			result = append(result, classes...)
			mutex.Unlock()
		}(course)
	}

	wg.Wait()

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *Plan) Schedules() (schedules [][]string, err error) {
	redisKey := "plan:" + strconv.FormatInt(p.ID, 10)
	cacheExists, err := db.Redis.Exists(redisKey).Result()

	if err != nil {
		return nil, err
	}

	// If cached schedule doesn't exist, regenerate the schedule
	if cacheExists == 0 {
		// Generate
		schedules, err = algorithm.Solve(db.UCRAPI, p.Term, p.CourseGroups)

		if err != nil {
			return nil, err
		}

		// Serialize
		serializedSchedules := make([]string, len(schedules))

		for i, sch := range schedules {
			str := strings.Join(sch, ",")
			serializedSchedules[i] = str
		}

		// Cache schedules in separate thread
		go func(redisKey string, serializedSchedules []string) {
			_, err = db.Redis.Del(redisKey).Result()

			if err != nil {
				//return nil, err
			}

			for i := 0; i < len(serializedSchedules); i += 500000 {
				var end int

				if i+500000 > len(serializedSchedules) {
					end = len(serializedSchedules)
				} else {
					end = i + 500000
				}

				_, err = db.Redis.LPush(redisKey, serializedSchedules[i:end]).Result()

				if err != nil {
					//return nil, err
				}
			}

			_, err = db.Redis.Expire(redisKey, 15*time.Minute).Result()

			if err != nil {
				//return nil, err
			}
		}(redisKey, serializedSchedules)

		return schedules, nil
	} else {
		serializedSchedules, err := db.Redis.LRange(redisKey, 0, -1).Result()

		if err != nil {
			return nil, err
		}

		schedules = make([][]string, len(serializedSchedules))

		for i, str := range serializedSchedules {
			schedules[i] = strings.Split(str, ",")
		}

		return schedules, nil
	}
}
