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

func devices(router *gin.Engine) {
	router.GET("/devices", devicesAll)
	router.GET("/devices/:id", devicesGetByID)
	router.POST("/devices", devicesCreate)
	router.PUT("/devices/:id", devicesUpdate)
	router.DELETE("/devices/:id", devicesDelete)
}

func devicesAll(c *gin.Context) {
	dc, err := models.NewDeviceController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	devices, err := dc.GetAll(ctx)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, devices)
}

func devicesGetByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		handleError(c, errors.NewServerError("No ID given/Invalid ID to get device", http.StatusNotAcceptable))
		return
	}

	dc, err := models.NewDeviceController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	device, err := dc.GetByID(ctx, id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, *device)
}

func devicesCreate(c *gin.Context) {
	var device models.Device
	dc, err := models.NewDeviceController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c.BindJSON(&device)

	if err := dc.Create(ctx, &device); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func devicesUpdate(c *gin.Context) {
	var device models.Device
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		handleError(c, errors.NewServerError("No ID given/Invalid ID to get device", http.StatusNotAcceptable))
		return
	}

	dc, err := models.NewDeviceController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c.BindJSON(&device)

	if err := dc.Update(ctx, id, &device); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{})
}

func devicesDelete(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		handleError(c, errors.NewServerError("No ID given/Invalid ID to get device", http.StatusNotAcceptable))
		return
	}

	dc, err := models.NewDeviceController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := dc.Delete(ctx, id); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
