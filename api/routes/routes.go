package routes

import (
	"github.com/gin-gonic/gin"
)

// Routes refines all the routes of the API
func Routes(router *gin.Engine) {
	devices(router)
}
