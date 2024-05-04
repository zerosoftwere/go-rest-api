package routes

import (
	"alhassan.link/rest-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	authenticated := server.Group("/")

	authenticated.Use(middlewares.Authenticate)
	authenticated.GET("/events", getEvents)
	authenticated.POST("/events", createEvent)
	authenticated.GET("/events/:id", getEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)

	authenticated.POST("/events/:id/register", registerForEvent)
	authenticated.DELETE("/events/:id/register", unregisterForEvent)
	authenticated.GET("/events/registrations", getEventRegistrations)

	server.POST("/signup", signup)
	server.POST("/login", login)
}
