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

func configurations(router *gin.Engine) {
	router.GET("/configurations", configurationsAll)
	router.GET("/configurations/by-id/:id", configurationsGetByID)
	router.GET("/configurations/by-user/:user", configurationsGetByUser)
	router.POST("/configurations", configurationsCreate)
	router.PUT("/configurations/update/:id", configurationsUpdate)
	router.DELETE("/configurations/:id", configurationsDelete)
	router.PUT("/configurations/set-off-update/:id", configurationsSetOffUpdate)
}

func configurationsAll(c *gin.Context) {
	cc, err := models.NewConfigurationController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	configurations, err := cc.GetAll(ctx)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, configurations)
}

func configurationsGetByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		handleError(c, errors.NewServerError("No ID given/Invalid ID to get device", http.StatusNotAcceptable))
		return
	}

	cc, err := models.NewConfigurationController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	configuration, err := cc.GetByID(ctx, id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, *configuration)
}

func configurationsGetByUser(c *gin.Context) {
	user, err := primitive.ObjectIDFromHex(c.Param("user"))
	if err != nil {
		handleError(c, errors.NewServerError("No user given/Invalid user to get device", http.StatusNotAcceptable))
		return
	}

	cc, err := models.NewConfigurationController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	configuration, err := cc.GetByUser(ctx, user)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, configuration)
}

func configurationsCreate(c *gin.Context) {
	var configuration models.Configuration
	cc, err := models.NewConfigurationController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c.BindJSON(&configuration)
	configuration.SystemStatus = "ACTIVE"
	configuration.RainPercentage = 80
	configuration.DefaultDuration = 10
	configuration.Update = true

	if err := cc.Create(ctx, &configuration); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func configurationsUpdate(c *gin.Context) {
	var configuration models.Configuration
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		handleError(c, errors.NewServerError("No ID given/Invalid ID to get configuration", http.StatusNotAcceptable))
		return
	}

	cc, err := models.NewConfigurationController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c.BindJSON(&configuration)

	if err := cc.Update(ctx, id, &configuration); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{})
}

func configurationsSetOffUpdate(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		handleError(c, errors.NewServerError("No ID given/Invalid ID to get configuration", http.StatusNotAcceptable))
		return
	}

	cc, err := models.NewConfigurationController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	configuration, err := cc.GetByID(ctx, id)
	if err != nil {
		handleError(c, err)
		return
	}

	configuration.Update = false

	if err := cc.Update(ctx, id, configuration); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{})
}

func configurationsDelete(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		handleError(c, errors.NewServerError("No ID given/Invalid ID to get configuration", http.StatusNotAcceptable))
		return
	}

	cc, err := models.NewConfigurationController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := cc.Delete(ctx, id); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
