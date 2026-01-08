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

// Login godoc
// @Summary Вход в аккаунт
// @Description Авторизация пользователя. Возвращает JWT токены (access для запросов + refresh для продления сессии). Токены сохраняются в HTTP-only cookies для безопасности
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Никнейм и пароль"
// @Success 200 {object} map[string]interface{} "Авторизация успешна! Добро пожаловать обратно"
// @Failure 400 {object} map[string]string "Не хватает данных (никнейм или пароль)"
// @Failure 401 {object} map[string]string "Неверный никнейм или пароль, попробуй еще раз"
// @Failure 429 {object} map[string]string "Слишком много попыток входа, подожди немного"
// @Failure 500 {object} map[string]string "Ошибка сервера при генерации токенов"
// @Router /login [post]
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

// RefreshAccessToken godoc
// @Summary Обновить токен
// @Description Обновляет access токен используя refresh токен. Вызывай этот эндпоинт когда access токен истек (обычно через 3 минуты). Refresh токен живет 30 дней
// @Tags Аутентификация
// @Produce json
// @Success 200 {object} map[string]string "Токен обновлен! Можешь продолжать работать"
// @Failure 401 {object} map[string]string "Refresh токен не найден, истек или невалидный. Нужно заново залогиниться"
// @Failure 500 {object} map[string]string "Ошибка генерации нового токена"
// @Router /refresh [post]
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

// Logout godoc
// @Summary Выход
// @Description Выход из аккаунта. Удаляет refresh токен из БД и чистит cookies. После этого все запросы будут отклонены, пока не залогинишься заново
// @Tags Аутентификация
// @Produce json
// @Success 200 {object} map[string]string "Успешный выход! До встречи"
// @Router /logout [post]
func Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")

	if refreshToken != "" {
		database.DB.Where("token = ?", refreshToken).Delete(&models.RefreshToken{})
	}

	utils.SetAuthCookie(c, "access_token", "", -1)
	utils.SetAuthCookie(c, "refresh_token", "", -1)

	c.JSON(http.StatusOK, gin.H{"message": "Успешный выход"})
}
