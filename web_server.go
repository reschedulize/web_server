package main

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/reschedulize/web_server/db"
	"github.com/reschedulize/web_server/middlewares"
	"github.com/reschedulize/web_server/routes_api"
	"github.com/reschedulize/web_server/routes_frontend"
)

func main() {
	db.ConnectMySQL()
	db.ConnectRedis()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Static("assets/", "build/assets/")
	r.LoadHTMLGlob("build/*.html")

	public := r.Group("/")
	public.Use(middlewares.Public)
	{
		public.GET("/", routes_frontend.GETLandingPage)

		public.GET("/login", routes_frontend.GETLogin)
		public.POST("/login", routes_frontend.POSTLogin)
	}

	private := r.Group("/")
	private.Use(middlewares.Private)
	{
		private.GET("/logout", routes_frontend.GETLogout)

		private.GET("/plans/", routes_frontend.GETViewAllPlans)
		private.GET("/plans/view/:id", routes_frontend.GETViewPlan)
		private.GET("/plans/new", routes_frontend.GETNewPlan)

		api := private.Group("/api")
		{
			api.GET("/user", routes_api.GETUser)

			api.GET("/plans", routes_api.GETPlanList)
			api.GET("/plans/:id", routes_api.GETPlan)
			api.POST("/plans", routes_api.POSTPlan)

			api.GET("/data/terms", routes_api.GETTerms)
			api.GET("/data/course_list/:term", routes_api.GETCourseList)
			api.GET("/data/course_stats/:term/:course", routes_api.GETCourseStats)
			api.GET("/data/classes/:term/:courses", routes_api.GETClasses)
		}
	}

	r.Run("0.0.0.0:8080")
}
