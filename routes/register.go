package routes

import (
	"net/http"
	"strconv"

	"alhassan.link/rest-api/models"
	"github.com/gin-gonic/gin"
)

func registerForEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse event id"})
		return
	}

	if models.HasEventRegistration(userId, eventId) {
		context.JSON(http.StatusBadRequest, gin.H{"message": "event already registered"})
		return
	}

	_, err = models.GetEvent(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "event not found"})
		return
	}
	registration := models.Registration{
		EventID: eventId,
		UserID:  userId,
	}
	err = registration.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not register event"})
		return
	}
	context.JSON(http.StatusCreated, registration)
}

func unregisterForEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse event id"})
		return
	}

	if !models.HasEventRegistration(userId, eventId) {
		context.JSON(http.StatusNotFound, gin.H{"message": "event registration not found"})
		return
	}

	err = models.DeleteEventRegistration(userId, eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not unregister for event"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"mesage": "event unregistered successfully"})
}

func getEventRegistrations(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventRegistations, err := models.GetEventRegistrations(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch event registrations"})
		return
	}
	context.JSON(http.StatusOK, eventRegistations)
}
