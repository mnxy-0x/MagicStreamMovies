package main

import (
	"log"

	"github.com/gin-gonic/gin"
	controller "github.com/mnxy-0x/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
)

func main() {
	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Starting the App...")
	})

	router.GET("/movies", controller.GetMovies())
	router.GET("/movie/:imdb_id", controller.GetOneMovie())
		// (imdb_id) should be the same name extracted by c.Param("imdb_id") in the handler [GetOneMovie()]

	router.POST("/addmovie", controller.AddMovie())
	router.POST("/register", controller.RegisterUser())
	router.POST("/login", controller.LoginUser())

	// every request handled by unique go_routine under-the-hood by (router.Run)
	if err := router.Run(":4000"); err != nil {
		log.Println("Failed to start server", err)
	}

}
