package middleware

import (
	"arizonagamesstore/backend/database"
	"arizonagamesstore/backend/models"
	"arizonagamesstore/backend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")

		if err == nil && accessToken != "" {
			claims, err := utils.ValidateAccessToken(accessToken)
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("nickname", claims.Nickname)
				c.Next()
				return
			}
		}

		refreshToken, err := c.Cookie("refresh_token")
		if err != nil || refreshToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Необходима авторизация"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateRefreshToken(refreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Необходима авторизация"})
			c.Abort()
			return
		}

		var storedToken models.RefreshToken
		result := database.DB.Where("token = ? AND account_id = ? AND expires_at > ?",
			refreshToken, claims.UserID, time.Now()).First(&storedToken)

		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Необходима авторизация"})
			c.Abort()
			return
		}

		newAccessToken, err := utils.GenerateAccessToken(claims.UserID, claims.Nickname)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
			c.Abort()
			return
		}

		utils.SetAuthCookie(c, "access_token", newAccessToken, 180)

		c.Set("user_id", claims.UserID)
		c.Set("nickname", claims.Nickname)
		c.Next()
	}
}
