package middleware

import (
	"academix/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthenticateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the token from the Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if present
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		// Verify the token
		token, claims, err := auth.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Ensure token is valid
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user data in context for next handlers
		c.Set("username", claims["username"].(string))
		c.Next() // Proceed to the next middleware or route handler
	}
}
