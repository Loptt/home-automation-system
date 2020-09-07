package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/Loptt/home-automation-system/api/auth"
	"github.com/Loptt/home-automation-system/api/errors"
	"github.com/Loptt/home-automation-system/api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func users(router *gin.Engine) {
	router.GET("/users", usersAll)
	router.GET("/users/:id", usersGetByID)
	router.POST("/users", usersCreate)
	router.PUT("/users/:id", usersUpdate)
	router.DELETE("/users/:id", usersDelete)
}

func usersAll(c *gin.Context) {
	uc, err := models.NewUserController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users, err := uc.GetAll(ctx)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func usersGetByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		handleError(c, errors.NewServerError("No ID given/Invalid ID to get user", http.StatusNotAcceptable))
		return
	}

	uc, err := models.NewUserController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := uc.GetByID(ctx, id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, *user)
}

func usersCreate(c *gin.Context) {
	var user models.User
	dc, err := models.NewUserController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c.BindJSON(&user)
	hashed, err := auth.HashPassword(user.Password)
	if err != nil {
		handleError(c, err)
		return
	}

	user.Password = hashed

	if err := dc.Create(ctx, &user); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func usersUpdate(c *gin.Context) {
	var user models.User
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		handleError(c, errors.NewServerError("No ID given/Invalid ID to get user", http.StatusNotAcceptable))
		return
	}

	dc, err := models.NewUserController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c.BindJSON(&user)
	hashed, err := auth.HashPassword(user.Password)
	if err != nil {
		handleError(c, err)
		return
	}

	user.Password = hashed

	if err := dc.Update(ctx, id, &user); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{})
}

func usersDelete(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		handleError(c, errors.NewServerError("No ID given/Invalid ID to get user", http.StatusNotAcceptable))
		return
	}

	dc, err := models.NewUserController()
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
