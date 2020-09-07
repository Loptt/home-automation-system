package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/Loptt/home-automation-system/api/auth"
	"github.com/Loptt/home-automation-system/api/models"
	"github.com/gin-gonic/gin"
)

func authentication(router *gin.Engine) {
	router.POST("/login", login)
}

func login(c *gin.Context) {
	var user models.User
	uc, err := models.NewUserController()
	if err != nil {
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c.BindJSON(&user)

	foundUser, err := uc.GetByUsername(ctx, user.Username)
	if err != nil {
		handleError(c, err)
		return
	}

	valid := auth.ValidatePassword(foundUser.Password, user.Password)

	c.JSON(http.StatusOK, gin.H{"valid": valid})
}
