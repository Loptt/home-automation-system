package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/Loptt/home-automation-system/api/errors"
	"github.com/Loptt/home-automation-system/api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func events(router *gin.Engine) {
	router.GET("/events", eventsAll)
	router.GET("/events/by-id/:id", eventsGetByID)
	router.GET("/events/by-device/:device", eventsGetByDevice)
	router.POST("/events", eventsCreate)
}

func eventsAll(c *gin.Context) {
	dc, err := models.NewEventController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	events, err := dc.GetAll(ctx)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, events)
}

func eventsGetByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		handleError(c, errors.NewServerError("No ID given/Invalid ID to get device", http.StatusNotAcceptable))
		return
	}

	dc, err := models.NewEventController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	event, err := dc.GetByID(ctx, id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, *event)
}

func eventsGetByDevice(c *gin.Context) {
	device, err := primitive.ObjectIDFromHex(c.Param("device"))
	if err != nil {
		handleError(c, errors.NewServerError("No device given/Invalid device to get device", http.StatusNotAcceptable))
		return
	}

	dc, err := models.NewEventController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	events, err := dc.GetByDevice(ctx, device)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, events)
}

func eventsCreate(c *gin.Context) {
	var event models.Event
	dc, err := models.NewEventController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c.BindJSON(&event)

	if err := dc.Create(ctx, &event); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}
