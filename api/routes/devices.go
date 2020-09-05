package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Loptt/home-automation-system/api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func devices(router *gin.Engine) {
	router.GET("/devices", devicesAll)
	router.GET("/devices/:id", devicesGetByID)
	router.POST("/devices", devicesCreate)
	router.DELETE("/devices/:id", devicesDelete)
}

func devicesAll(c *gin.Context) {
	dc, err := models.NewDeviceController()
	if err != nil {
		log.Println("Error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	devices, err := dc.GetAll(ctx)
	if err != nil {
		log.Println("Error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, devices)
}

func devicesGetByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Println("Error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	dc, err := models.NewDeviceController()
	if err != nil {
		log.Println("Error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	device, err := dc.GetByID(ctx, id)
	if err != nil {
		log.Println("Error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, *device)
}

func devicesCreate(c *gin.Context) {
	var device models.Device
	dc, err := models.NewDeviceController()
	if err != nil {
		log.Println("Error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c.BindJSON(&device)

	if err := dc.Create(ctx, &device); err != nil {
		log.Println("Error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func devicesDelete(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Println("Error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	dc, err := models.NewDeviceController()
	if err != nil {
		log.Println("Error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := dc.Delete(ctx, id); err != nil {
		log.Println("Error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
