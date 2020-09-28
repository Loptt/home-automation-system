package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Routes refines all the routes of the API
func Routes(router *gin.Engine) {
	router.GET("/", defaultRoute)
	users(router)
	devices(router)
	events(router)
	configurations(router)
	authentication(router)
}

func defaultRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Home Automation System API"})
}
