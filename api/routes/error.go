package routes

import (
	"log"
	"net/http"

	"github.com/Loptt/home-automation-system/api/errors"
	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, err error) {
	log.Println("Error: " + err.Error())

	switch err.(type) {
	case *errors.ServerError:
		c.AbortWithStatusJSON(err.(*errors.ServerError).Code(), gin.H{"message": err.(*errors.ServerError).Message()})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
	}
}
