package main

import (
	"log"

	"github.com/Loptt/home-automation-system/api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	routes.Routes(router)
	log.Fatal(router.Run(":4747"))
}
