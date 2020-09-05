package main

import (
	"log"

	"github.com/Loptt/home-automation-system/api/config"
	"github.com/Loptt/home-automation-system/api/db"
	"github.com/Loptt/home-automation-system/api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()
	router := gin.Default()

	port, err := config.Port()
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}

	dburl, err := config.DatabaseURL()
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}

	if err := db.Connect(dburl); err != nil {
		log.Fatal("Error: " + err.Error())
	}

	routes.Routes(router)
	log.Fatal(router.Run(port))
}
