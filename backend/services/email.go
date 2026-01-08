package services

import (
	"arizonagamesstore/backend/database"
	"arizonagamesstore/backend/models"
	"arizonagamesstore/backend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type VerifyEmailRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Code    string `json:"code" binding:"required"`
	ClientIP string `json:"client_ip"`
}

type ResendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyEmail godoc
// @Summary Подтвердить email
// @Description Подтверждает email пользователя после регистрации. Нужно ввести код который пришел на почту. После подтверждения сразу логинит пользователя и выдает токены
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param request body VerifyEmailRequest true "Email и код подтверждения"
// @Success 200 {object} map[string]interface{} "Email подтвержден! Добро пожаловать в Arizona Games Store"
// @Failure 400 {object} map[string]string "Неверный код или он уже истек (коды живут 10 минут)"
// @Failure 403 {object} map[string]string "Регистрация заблокирована (превышен лимит с одного IP)"
// @Failure 500 {object} map[string]string "Ошибка создания аккаунта или генерации токенов"
// @Router /verify-email [post]
func VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат. Укажите email и код."})
		return
	}

	if blockedCookie, err := c.Cookie("reg_blocked"); err == nil && blockedCookie != "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Регистрация заблокирована. Обратитесь в поддержку."})
		return
	}

	var verification models.EmailVerification
	if err := database.DB.Where("email = ? AND code = ? AND expires_at > ?", req.Email, req.Code, time.Now()).First(&verification).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный код подтверждения или срок его действия истёк."})
		return
	}

	clientIP := req.ClientIP
	if clientIP == "" {
		clientIP = c.ClientIP()
	}

	if clientIP == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось определить IP адрес. Попробуйте позже."})
		return
	}

	if clientIP != "::1" && clientIP != "127.0.0.1" {
		count, err := CountAccountsByIP(clientIP)
		if err == nil && count >= 3 {
			c.SetCookie(
				"reg_blocked",
				"1",
				365*24*60*60,
				"/",
				"",
				false,
				true,
			)
			c.JSON(http.StatusForbidden, gin.H{"error": "Достигнут лимит регистраций с вашего IP адреса (максимум 3 аккаунта)."})
			return
		}
	}

	if err := CreateAccountWithEmail(verification.Nickname, verification.Email, verification.PasswordHash, clientIP, clientIP, true); err != nil {
		if utils.IsDuplicateKeyError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email уже зарегистрирован или используется"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании аккаунта"})
		return
	}

	account, err := GetUserByNickname(verification.Nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных аккаунта"})
		return
	}

	database.DB.Where("email = ?", req.Email).Delete(&models.EmailVerification{})

	accessToken, err := utils.GenerateAccessToken(account.ID, account.Nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании access токена"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(account.ID, account.Nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании refresh токена"})
		return
	}

	tokenRecord := models.RefreshToken{
		AccountID: account.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := database.DB.Create(&tokenRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении refresh токена"})
		return
	}

	utils.SetAuthCookie(c, "access_token", accessToken, 180)
	utils.SetAuthCookie(c, "refresh_token", refreshToken, 30*24*60*60)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Email успешно подтверждён",
		"nickname": account.Nickname,
		"user_id":  account.ID,
	})
}

// ResendVerificationCode godoc
// @Summary Отправить код повторно
// @Description Отправляет новый код подтверждения на email. Нужно если предыдущий код истек (они живут 10 минут) или потерялся. Генерирует новый код и сразу отправляет на почту
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param request body ResendCodeRequest true "Email для отправки нового кода"
// @Success 200 {object} map[string]string "Новый код отправлен! Проверь почту (и спам тоже)"
// @Failure 400 {object} map[string]string "Email не найден или уже подтвержден"
// @Failure 500 {object} map[string]string "Ошибка отправки email"
// @Router /resend-code [post]
func ResendVerificationCode(c *gin.Context) {
	var req ResendCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат. Укажите email."})
		return
	}

	var verification models.EmailVerification
	if err := database.DB.Where("email = ?", req.Email).First(&verification).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email не найден или уже подтверждён"})
		return
	}

	verification.Code = utils.GenerateVerificationCode()
	verification.ExpiresAt = time.Now().Add(10 * time.Minute)

	if err := database.DB.Save(&verification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении кода подтверждения"})
		return
	}

	if err := utils.SendVerificationEmail(req.Email, verification.Code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при отправке повторного email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Новый код подтверждения отправлен"})
}
