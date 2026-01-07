package services

import (
	"arizonagamesstore/backend/database"
	"arizonagamesstore/backend/models"
	"arizonagamesstore/backend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Nickname       string `json:"nickname" binding:"required"`
	Password       string `json:"password" binding:"required"`
	RecaptchaToken string `json:"recaptcha_token"`
	ClientIP       string `json:"client_ip"`
}

func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	if req.RecaptchaToken != "" {
		valid, score, err := utils.VerifyRecaptcha(req.RecaptchaToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки reCAPTCHA"})
			return
		}
		if !valid || score < 0.5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Проверка reCAPTCHA не пройдена"})
			return
		}
	}

	var account models.Account
	result := database.DB.Where("nickname = ?", req.Nickname).First(&account)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный никнейм или пароль"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный никнейм или пароль"})
		return
	}

	accessToken, err := utils.GenerateAccessToken(account.ID, account.Nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(account.ID, account.Nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	tokenRecord := models.RefreshToken{
		AccountID: account.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := database.DB.Create(&tokenRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения токена"})
		return
	}

	utils.SetAuthCookie(c, "access_token", accessToken, 180)
	utils.SetAuthCookie(c, "refresh_token", refreshToken, 30*24*60*60)

	clientIP := req.ClientIP
	if clientIP == "" {
		clientIP = c.ClientIP()
	}

	if clientIP != "" {
		if err := UpdateLastIP(account.Nickname, clientIP); err != nil {
			println("Предупреждение: не удалось обновить last_ip:", err.Error())
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Успешная авторизация",
		"nickname": account.Nickname,
		"user_id":  account.ID,
	})
}

func RefreshAccessToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh токен не найден"})
		return
	}

	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидный refresh токен"})
		return
	}

	var storedToken models.RefreshToken
	result := database.DB.Where("token = ? AND account_id = ? AND expires_at > ?",
		refreshToken, claims.UserID, time.Now()).First(&storedToken)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh токен не найден или истек"})
		return
	}

	newAccessToken, err := utils.GenerateAccessToken(claims.UserID, claims.Nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	utils.SetAuthCookie(c, "access_token", newAccessToken, 180)

	c.JSON(http.StatusOK, gin.H{"message": "Токен обновлен"})
}

func Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")

	if refreshToken != "" {
		database.DB.Where("token = ?", refreshToken).Delete(&models.RefreshToken{})
	}

	utils.SetAuthCookie(c, "access_token", "", -1)
	utils.SetAuthCookie(c, "refresh_token", "", -1)

	c.JSON(http.StatusOK, gin.H{"message": "Успешный выход"})
}
