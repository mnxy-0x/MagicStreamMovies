package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mnxy-0x/MagicStreamMovies/Server/MagicStreamMoviesServer/utils"
)

// middleware function that protects your API routes by ensuring only requests with valid JWT access tokens can proceed
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.GetAccessToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token Provided"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "debug": err.Error()})
			c.Abort()
			return
		}

		// Store User Info in Gin context for later use
		// by this, User data is available to all the route handlers
		c.Set("userId", claims.UserId)
		c.Set("role", claims.Role)
		c.Next() // Without this, the request would stop here (even if authenticated)
	}
}
