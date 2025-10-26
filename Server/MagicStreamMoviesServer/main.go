package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mnxy-0x/MagicStreamMovies/Server/MagicStreamMoviesServer/routes"
)

func main() {
	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Starting the App...")
	})

	routes.SetupRoutes(router)

	// every request handled by unique go_routine under-the-hood by (router.Run)
	if err := router.Run(":4000"); err != nil {
		log.Println("Failed to start server", err)
	}

}
