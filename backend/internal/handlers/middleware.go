package handlers

import (
	"net/http"

	"chat-system/backend/internal/auth"
	"chat-system/backend/internal/config"

	"github.com/gin-gonic/gin"
)

func WSAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "token requerido como query param"})
			return
		}

		claims, err := auth.VerifyJWT(token, cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "El token que diste es inválido o expirado"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("pin", claims.Pin)
		c.Next()
	}
}
