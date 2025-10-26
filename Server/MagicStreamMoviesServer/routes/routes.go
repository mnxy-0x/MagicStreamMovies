package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/mnxy-0x/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
	"github.com/mnxy-0x/MagicStreamMovies/Server/MagicStreamMoviesServer/middlewares"
)

func SetupRoutes(router *gin.Engine) {

	public := router.Group("/")
	{
		public.GET("/movies", controller.GetMovies())
		public.POST("/register", controller.RegisterUser())
		public.POST("/login", controller.LoginUser())
	}

	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/movie/:imdb_id", controller.GetOneMovie()) // (imdb_id) should equal name extracted by c.Param("imdb_id") in the handler [GetOneMovie()]

		protected.POST("/addmovie", controller.AddMovie())
	}
}
