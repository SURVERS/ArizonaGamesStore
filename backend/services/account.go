package services

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"

	"arizonagamesstore/backend/database"
	"arizonagamesstore/backend/models"
	"arizonagamesstore/backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Nickname       string `json:"nickname" binding:"required"`
	Password       string `json:"password" binding:"required,min=6"`
	Email          string `json:"email" binding:"required,email"`
	RecaptchaToken string `json:"recaptcha_token"`
}

var sqlKeywords = []string{
	"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER",
	"EXEC", "EXECUTE", "UNION", "JOIN", "WHERE", "FROM", "TABLE",
	"DATABASE", "SCHEMA", "OR", "AND", "--", "/*", "*/", "xp_",
	"sp_", "CAST", "CONVERT", "CHAR", "NCHAR", "VARCHAR", "NVARCHAR",
	"TRUNCATE", "GRANT", "REVOKE", "INFORMATION_SCHEMA", "SYSOBJECTS",
	"SYSCOLUMNS", "0x", "WAITFOR", "DELAY", "BENCHMARK", "SLEEP",
}

var dangerousPatterns = []string{
	"<script", "</script>", "<iframe", "</iframe>", "javascript:",
	"onerror=", "onclick=", "onload=", "<img", "<svg", "<object",
	"<embed", "<link", "<style", "</style>", "eval(", "alert(",
	"prompt(", "confirm(", "document.", "window.", "<meta",
	"onmouseover=", "onfocus=", "onblur=", "<base", "<form",
	"expression(", "vbscript:", "data:", "<applet", "<bgsound",
}

func validateNickname(nickname string) (bool, string) {
	if strings.TrimSpace(nickname) == "" {
		return false, "Никнейм не может быть пустым"
	}

	if len(nickname) < 3 {
		return false, "Никнейм должен содержать минимум 3 символа"
	}
	if len(nickname) > 20 {
		return false, "Никнейм не должен превышать 20 символов"
	}

	trimmed := strings.TrimSpace(nickname)
	if trimmed != nickname {
		return false, "Никнейм не должен содержать пробелы в начале или конце"
	}
	if strings.Contains(nickname, " ") {
		return false, "Никнейм не должен содержать пробелы"
	}

	nicknameRegex := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*[a-zA-Z0-9_]$`)
	if len(nickname) == 1 {
		nicknameRegex = regexp.MustCompile(`^[a-zA-Z]$`)
	}
	if !nicknameRegex.MatchString(nickname) {
		return false, "Никнейм может содержать только английские буквы, цифры и подчёркивание (_). Должен начинаться с буквы"
	}

	if strings.Contains(nickname, "--") {
		return false, "Никнейм не должен содержать несколько дефисов подряд"
	}

	for _, char := range nickname {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' {
			return false, "Никнейм содержит недопустимые символы. Разрешены только английские буквы, цифры и дефис (_)"
		}
	}

	for _, char := range nickname {
		if char > 127 {
			return false, "Никнейм может содержать только английские буквы"
		}
		if unicode.IsLetter(char) {
			if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')) {
				return false, "Никнейм может содержать только английские буквы"
			}
		}
	}

	upperNickname := strings.ToUpper(nickname)
	for _, keyword := range sqlKeywords {
		if strings.Contains(upperNickname, keyword) {
			return false, "Никнейм содержит недопустимые символы или слова"
		}
	}

	lowerNickname := strings.ToLower(nickname)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerNickname, strings.ToLower(pattern)) {
			return false, "Никнейм содержит недопустимые символы"
		}
	}

	suspiciousPatterns := []string{"../", "..\\", "./", ".\\", "//", "\\\\"}
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(nickname, pattern) {
			return false, "Никнейм содержит недопустимые последовательности символов"
		}
	}

	return true, ""
}

func validatePassword(password string) (bool, string) {
	if strings.TrimSpace(password) == "" {
		return false, "Пароль не может быть пустым"
	}

	if len(password) < 6 {
		return false, "Пароль должен содержать минимум 6 символов"
	}
	if len(password) > 100 {
		return false, "Пароль не должен превышать 100 символов"
	}

	upperPassword := strings.ToUpper(password)
	for _, keyword := range sqlKeywords {
		if strings.Contains(upperPassword, keyword) {
			return false, "Пароль содержит недопустимые символы или конструкции"
		}
	}

	lowerPassword := strings.ToLower(password)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerPassword, strings.ToLower(pattern)) {
			return false, "Пароль содержит недопустимые символы или конструкции"
		}
	}

	if strings.Contains(password, "<") || strings.Contains(password, ">") {
		return false, "Пароль не должен содержать символы < или >"
	}

	if strings.Contains(password, "\x00") {
		return false, "Пароль содержит недопустимые символы"
	}
	for _, char := range password {
		if char < 32 && char != 9 && char != 10 && char != 13 {
			return false, "Пароль содержит недопустимые управляющие символы"
		}
	}

	suspiciousPatterns := []string{
		"\\'", "\\\"", "%00", "../", "..\\", "${", "#{",
		"<!--", "-->", "/*", "*/", "@@", "0x", "--",
	}
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(password, pattern) {
			return false, "Пароль содержит недопустимые последовательности символов"
		}
	}

	return true, ""
}

func RegisterAccount(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных. Проверьте правильность заполнения полей"})
		return
	}

	if blockedCookie, err := c.Cookie("reg_blocked"); err == nil && blockedCookie != "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Регистрация заблокирована. Обратитесь в поддержку."})
		return
	}

	clientIP := c.ClientIP()
	if clientIP != "" && clientIP != "::1" && clientIP != "127.0.0.1" {
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

	if valid, msg := validateNickname(req.Nickname); !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	if valid, msg := validatePassword(req.Password); !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	var existingAccount models.Account
	err := database.DB.Where("email = ?", req.Email).First(&existingAccount).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email уже используется"})
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверке email"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при регистрации аккаунта. Код ошибки: 001. Попробуйте позже."})
		return
	}

	var existingNickname models.Account
	err = database.DB.Where("nickname = ?", req.Nickname).First(&existingNickname).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Аккаунт с таким никнеймом уже существует"})
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверке никнейма"})
		return
	}

	database.DB.Where("email = ? OR nickname = ?", req.Email, req.Nickname).Delete(&models.EmailVerification{})

	code := utils.GenerateVerificationCode()
	verification := models.EmailVerification{
		Email:        req.Email,
		Nickname:     req.Nickname,
		PasswordHash: string(hashedPassword),
		Code:         code,
		ExpiresAt:    time.Now().Add(10 * time.Minute),
	}

	if err := database.DB.Create(&verification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании кода подтверждения"})
		return
	}

	if err := utils.SendVerificationEmail(req.Email, code); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message":         req.Nickname + " успешно зарегистрирован!",
			"nickname":        req.Nickname,
			"email_sent":      false,
			"email_error":     "Не удалось отправить письмо с кодом подтверждения",
			"email_verified":  false,
			"requires_verify": true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         req.Nickname + " успешно зарегистрирован!",
		"nickname":        req.Nickname,
		"email_sent":      true,
		"email_verified":  false,
		"requires_verify": true,
	})
}

func GetUserByNickname(nickname string) (*models.Account, error) {
	var account models.Account
	err := database.DB.Where("nickname = ?", nickname).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func UpdateLastSeen(nickname string) error {
	return database.DB.Model(&models.Account{}).
		Where("nickname = ?", nickname).
		Update("last_seen_at", time.Now()).Error
}

func UpdateUserAvatar(nickname string, avatarURL string) error {
	return database.DB.Model(&models.Account{}).
		Where("nickname = ?", nickname).
		Update("avatar", avatarURL).Error
}

// UpdateUserNickname обновляет никнейм пользователя и все его объявления
func UpdateUserNickname(oldNickname string, newNickname string) error {
	now := time.Now()

	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&models.Account{}).
		Where("nickname = ?", oldNickname).
		Updates(map[string]interface{}{
			"nickname":             newNickname,
			"last_nickname_change": &now,
			"last_settings_change": &now,
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func UpdateUserEmail(nickname string, newEmail string) error {
	now := time.Now()
	return database.DB.Model(&models.Account{}).
		Where("nickname = ?", nickname).
		Updates(map[string]interface{}{
			"email":                newEmail,
			"last_email_change":    &now,
			"last_settings_change": &now,
		}).Error
}

func UpdateUserPassword(nickname string, newPasswordHash string) error {
	now := time.Now()
	return database.DB.Model(&models.Account{}).
		Where("nickname = ?", nickname).
		Updates(map[string]interface{}{
			"password_hash":        newPasswordHash,
			"last_settings_change": &now,
		}).Error
}

func UpdateUserTheme(nickname string, theme string) error {
	return database.DB.Model(&models.Account{}).
		Where("nickname = ?", nickname).
		Update("theme", theme).Error
}

func UpdateUserDescription(nickname string, description string) error {
	now := time.Now()
	return database.DB.Model(&models.Account{}).
		Where("nickname = ?", nickname).
		Updates(map[string]interface{}{
			"user_description":     description,
			"last_settings_change": &now,
		}).Error
}

func CheckNicknameExists(nickname string) (bool, error) {
	var account models.Account
	err := database.DB.Where("nickname = ?", nickname).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func CheckEmailExists(email string) (bool, error) {
	var account models.Account
	err := database.DB.Where("email = ?", email).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func CreateAccount(nickname string, passwordHash string, regIP string, lastIP string) error {
	account := models.Account{
		Nickname:      nickname,
		PasswordHash:  passwordHash,
		Theme:         "dark",
		RegIP:         regIP,
		LastIP:        lastIP,
		UserRole:      "User",
		EmailVerified: false,
		Avatar:        "https://storage.yandexcloud.net/fotora.ru/uploads/2b0c131e8cfe54b1.jpeg",
	}

	result := database.DB.Create(&account)
	return result.Error
}

func CreateAccountWithEmail(nickname string, email string, passwordHash string, regIP string, lastIP string, emailVerified bool) error {
	account := models.Account{
		Nickname:      nickname,
		Email:         email,
		EmailVerified: emailVerified,
		PasswordHash:  passwordHash,
		Theme:         "dark",
		RegIP:         regIP,
		LastIP:        lastIP,
		UserRole:      "User",
		Avatar:        "https://storage.yandexcloud.net/fotora.ru/uploads/2b0c131e8cfe54b1.jpeg",
	}

	result := database.DB.Create(&account)
	return result.Error
}

func UpdateProfileBackground(nickname string, backgroundURL string) error {
	return database.DB.Model(&models.Account{}).
		Where("nickname = ?", nickname).
		Update("background_avatar_profile", backgroundURL).Error
}

func DeleteProfileBackground(nickname string) error {
	return database.DB.Model(&models.Account{}).
		Where("nickname = ?", nickname).
		Update("background_avatar_profile", "").Error
}

func CountAccountsByIP(ip string) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Account{}).
		Where("reg_ip = ?", ip).
		Count(&count).Error
	return count, err
}

func UpdateLastIP(nickname string, ip string) error {
	return database.DB.Model(&models.Account{}).
		Where("nickname = ?", nickname).
		Update("last_ip", ip).Error
}

func UpdateTelegram(nickname string, telegram string) error {
	now := time.Now()
	return database.DB.Model(&models.Account{}).
		Where("nickname = ?", nickname).
		Updates(map[string]interface{}{
			"telegram":             telegram,
			"last_settings_change": &now,
		}).Error
}
