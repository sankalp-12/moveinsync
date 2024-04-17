package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sankalp-12/moveinsync/user-service/utils"
)

// It is a auth-middleware to authenticate requests using JWT token for admin enpoints
func Validate(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get JWT token from the request header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Error().Msg("Internal server error: Authorization header is missing")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Extract the token from the Authorization header
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.Error().Msg("Internal server error: Invalid authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}
		tokenString := tokenParts[1]

		// Verify the JWT signature
		valid, err := utils.VerifyJWT(tokenString)
		if err != nil {
			logger.Error().Msg("Internal server error: Failed to verify JWT signature")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify JWT signature"})
			c.Abort()
			return
		}

		if !valid {
			logger.Error().Msg("Internal server error: Invalid JWT token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token"})
			c.Abort()
			return
		}

		// Token is valid, proceed with the next middleware or handler
		logger.Info().Msg("User request validated successfully")
		c.Next()
	}
}
