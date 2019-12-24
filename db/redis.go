package db

import (
	"github.com/go-redis/redis"
	"github.com/reschedulize/school_course_data"
)

var Redis *redis.Client
var UCRAPI *school_course_data.UCRAPI

func ConnectRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	err := Redis.Ping().Err()

	if err != nil {
		panic(err)
	}

	UCRAPI = school_course_data.NewUCRAPI(Redis)
}
