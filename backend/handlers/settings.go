package handlers

import (
	"arizonagamesstore/backend/services"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UpdateNickname godoc
// @Summary Изменить никнейм
// @Description Изменяет никнейм пользователя. Макс. 20 символов
// @Tags Профиль
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Новый никнейм" example(nickname="NewNick123")
// @Success 200 {object} map[string]string "Никнейм обновлен!"
// @Failure 400 {object} map[string]string "Никнейм слишком длинный или пустой"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 409 {object} map[string]string "Такой никнейм уже занят"
// @Failure 500 {object} map[string]string "Ошибка обновления"
// @Router /profile/update-nickname [put]
func UpdateNickname(c *gin.Context) {
	nickname, exists := c.Get("nickname")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	var req struct {
		Nickname string `json:"nickname" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Получаем текущего пользователя
	user, err := services.GetUserByNickname(nickname.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных пользователя"})
		return
	}

	// Проверка cooldown - 1 неделя для никнейма
	if user.LastNicknameChange != nil {
		timeSinceLastChange := time.Since(*user.LastNicknameChange)
		if timeSinceLastChange < 7*24*time.Hour {
			remainingTime := 7*24*time.Hour - timeSinceLastChange
			days := int(remainingTime.Hours() / 24)
			hours := int(remainingTime.Hours()) % 24
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("Вы можете изменить никнейм только раз в неделю. Осталось: %d дн. %d ч.", days, hours),
			})
			return
		}
	}

	// Валидация нового никнейма
	if len(req.Nickname) < 3 || len(req.Nickname) > 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Никнейм должен содержать от 3 до 20 символов"})
		return
	}

	nicknameRegex := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	if !nicknameRegex.MatchString(req.Nickname) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Никнейм может содержать только английские буквы, цифры и подчёркивание"})
		return
	}

	// Проверка уникальности никнейма через сервисный слой
	nicknameExists, err := services.CheckNicknameExists(req.Nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки никнейма"})
		return
	}
	if nicknameExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Никнейм уже занят"})
		return
	}

	// Обновляем никнейм
	if err := services.UpdateUserNickname(nickname.(string), req.Nickname); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления никнейма"})
		return
	}

	if err := services.UpdateNickNameAds(nickname.(string), req.Nickname); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления никнейма в ads"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Никнейм успешно обновлен"})
}

// UpdateEmail godoc
// @Summary Изменить email
// @Description Изменяет email пользователя. После изменения нужно будет заново подтвердить email
// @Tags Профиль
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Новый email" example(email="new@arizona.rp")
// @Success 200 {object} map[string]string "Email обновлен! Проверь почту для подтверждения"
// @Failure 400 {object} map[string]string "Некорректный email"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 409 {object} map[string]string "Такой email уже используется"
// @Failure 500 {object} map[string]string "Ошибка обновления"
// @Router /profile/update-email [put]
func UpdateEmail(c *gin.Context) {
	nickname, exists := c.Get("nickname")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат email"})
		return
	}

	// Получаем текущего пользователя
	user, err := services.GetUserByNickname(nickname.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных пользователя"})
		return
	}

	// Проверка cooldown - 1 неделя для email
	if user.LastEmailChange != nil {
		timeSinceLastChange := time.Since(*user.LastEmailChange)
		if timeSinceLastChange < 7*24*time.Hour {
			remainingTime := 7*24*time.Hour - timeSinceLastChange
			days := int(remainingTime.Hours() / 24)
			hours := int(remainingTime.Hours()) % 24
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("Вы можете изменить email только раз в неделю. Осталось: %d дн. %d ч.", days, hours),
			})
			return
		}
	}

	// Проверка уникальности email через сервисный слой
	emailExists, err := services.CheckEmailExists(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки email"})
		return
	}
	if emailExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email уже используется"})
		return
	}

	// Обновляем email
	if err := services.UpdateUserEmail(nickname.(string), req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email успешно обновлен"})
}

// UpdatePassword godoc
// @Summary Изменить пароль
// @Description Изменяет пароль пользователя. Нужно ввести старый пароль для подтверждения
// @Tags Профиль
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Старый и новый пароли" example(old_password="OldPass123!" new_password="NewPass123!")
// @Success 200 {object} map[string]string "Пароль обновлен!"
// @Failure 400 {object} map[string]string "Новый пароль слишком короткий (мин. 8 символов)"
// @Failure 401 {object} map[string]string "Старый пароль неверный"
// @Failure 500 {object} map[string]string "Ошибка обновления"
// @Router /profile/update-password [put]
func UpdatePassword(c *gin.Context) {
	nickname, exists := c.Get("nickname")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	var req struct {
		OldPassword     string `json:"oldPassword" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required,min=6"`
		ConfirmPassword string `json:"confirmPassword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Проверка совпадения паролей
	if req.NewPassword != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пароли не совпадают"})
		return
	}

	// Получаем текущего пользователя
	user, err := services.GetUserByNickname(nickname.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных пользователя"})
		return
	}

	// Проверка cooldown - 1 минута для пароля
	if user.LastSettingsChange != nil {
		timeSinceLastChange := time.Since(*user.LastSettingsChange)
		if timeSinceLastChange < 1*time.Minute {
			remainingTime := 1*time.Minute - timeSinceLastChange
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("Подождите %d секунд перед следующим изменением", int(remainingTime.Seconds())),
			})
			return
		}
	}

	// Проверяем старый пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный старый пароль"})
		return
	}

	// Хешируем новый пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании нового пароля"})
		return
	}

	// Обновляем пароль
	if err := services.UpdateUserPassword(nickname.(string), string(hashedPassword)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления пароля"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пароль успешно обновлен"})
}

// UpdateTheme godoc
// @Summary Изменить тему
// @Description Меняет тему оформления (светлая/темная)
// @Tags Профиль
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Тема" example(theme="dark")
// @Success 200 {object} map[string]string "Тема обновлена!"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 500 {object} map[string]string "Ошибка обновления"
// @Router /profile/update-theme [put]
func UpdateTheme(c *gin.Context) {
	nickname, exists := c.Get("nickname")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	var req struct {
		Theme string `json:"theme" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Валидация темы
	if req.Theme != "dark" && req.Theme != "light" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Допустимые значения: dark или light"})
		return
	}

	// Обновляем тему
	if err := services.UpdateUserTheme(nickname.(string), req.Theme); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления темы"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Тема успешно обновлена"})
}

// UpdateDescription godoc
// @Summary Изменить описание
// @Description Обновляет описание профиля. Макс. 500 символов
// @Tags Профиль
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Описание" example(description="Активный продавец на Arizona RP")
// @Success 200 {object} map[string]string "Описание обновлено!"
// @Failure 400 {object} map[string]string "Описание слишком длинное"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 500 {object} map[string]string "Ошибка обновления"
// @Router /profile/update-description [put]
func UpdateDescription(c *gin.Context) {
	nickname, exists := c.Get("nickname")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	var req struct {
		Description string `json:"description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Валидация длины описания
	if len(req.Description) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Описание должно содержать минимум 3 символа"})
		return
	}

	if len(req.Description) > 200 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Описание не должно превышать 200 символов"})
		return
	}

	// Санитизация: запрет HTML тегов, скриптов и SQL
	sanitized := sanitizeDescription(req.Description)

	// Проверка на опасные символы и паттерны
	if containsDangerousPatterns(sanitized) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Описание содержит запрещенные символы или паттерны"})
		return
	}

	// Получаем текущего пользователя для проверки cooldown
	user, err := services.GetUserByNickname(nickname.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных пользователя"})
		return
	}

	// Проверка cooldown - 1 минута
	if user.LastSettingsChange != nil {
		timeSinceLastChange := time.Since(*user.LastSettingsChange)
		if timeSinceLastChange < 1*time.Minute {
			remainingTime := 1*time.Minute - timeSinceLastChange
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("Подождите %d секунд перед следующим изменением", int(remainingTime.Seconds())),
			})
			return
		}
	}

	// Обновляем описание
	if err := services.UpdateUserDescription(nickname.(string), sanitized); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления описания"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Описание успешно обновлено"})
}

// sanitizeDescription очищает описание от опасного содержимого
func sanitizeDescription(description string) string {
	// Удаляем HTML теги
	re := regexp.MustCompile(`<[^>]*>`)
	cleaned := re.ReplaceAllString(description, "")

	// Удаляем JavaScript
	re = regexp.MustCompile(`(?i)<script.*?</script>`)
	cleaned = re.ReplaceAllString(cleaned, "")

	// Удаляем SQL ключевые слова (базовая защита)
	sqlKeywords := []string{"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER", "UNION", "EXEC", "EXECUTE"}
	for _, keyword := range sqlKeywords {
		re = regexp.MustCompile(`(?i)\b` + keyword + `\b`)
		cleaned = re.ReplaceAllString(cleaned, "")
	}

	return strings.TrimSpace(cleaned)
}

// containsDangerousPatterns проверяет наличие опасных паттернов
func containsDangerousPatterns(text string) bool {
	dangerousPatterns := []string{
		`<script`,
		`javascript:`,
		`onerror=`,
		`onload=`,
		`onclick=`,
		`eval\(`,
		`expression\(`,
		`vbscript:`,
		`<iframe`,
		`<object`,
		`<embed`,
		`1=1`,
		`' OR '`,
		`" OR "`,
		`--`,
		`;--`,
		`/*`,
		`*/`,
		`xp_`,
		`sp_`,
	}

	lowerText := strings.ToLower(text)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerText, strings.ToLower(pattern)) {
			return true
		}
	}

	return false
}

// UpdateProfileAvatar godoc
// @Summary Изменить аватар
// @Description Загружает новый аватар на S3. Макс. размер 5MB
// @Tags Профиль
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "Изображение аватара"
// @Success 200 {object} map[string]string "URL нового аватара"
// @Failure 400 {object} map[string]string "Файл не загружен или слишком большой"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 500 {object} map[string]string "Ошибка загрузки"
// @Router /profile/update-avatar [post]
func UpdateProfileAvatar(c *gin.Context) {
	nickname, exists := c.Get("nickname")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении файла"})
		return
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимый тип файла. Разрешены только JPEG, PNG и GIF."})
		return
	}

	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Размер файла превышает 5 МБ."})
		return
	}

	uniqueFileName := generateUniqueFileName(file.Filename)
	imagePath := fmt.Sprintf("avatars/%s/%s", nickname.(string), uniqueFileName)

	tempDir := "./temp_uploads"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания временной директории"})
		return
	}

	tempFilePath := filepath.Join(tempDir, uniqueFileName)
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения файла"})
		return
	}

	s3URL, err := UploadFileToS3(tempFilePath, imagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки файла в хранилище"})
		return
	}

	defer os.Remove(tempFilePath)

	user, err := services.GetUserByNickname(nickname.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных пользователя"})
		return
	}

	if user.Avatar != "" {
		oldKey := extractS3KeyFromURL(user.Avatar)
		if oldKey != "" {
			if err := DeleteFileFromS3(oldKey); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления старого аватара"})
				return
			}
		}
	}

	if err := services.UpdateUserAvatar(nickname.(string), s3URL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления данных пользователя"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Аватар успешно обновлен", "avatar_url": s3URL})
}
