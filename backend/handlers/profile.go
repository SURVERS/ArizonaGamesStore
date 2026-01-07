package handlers

import (
	"arizonagamesstore/backend/services"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func UpdateProfileBackground(c *gin.Context) {
	// Получаем nickname из контекста
	nickname, exists := c.Get("nickname")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	// Получаем текущего пользователя из БД
	user, err := services.GetUserByNickname(nickname.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных пользователя"})
		return
	}

	// Получение файла изображения
	file, err := c.FormFile("background")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Изображение обязательно"})
		return
	}

	// Проверка размера файла (20MB)
	if file.Size > 20*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Размер изображения не должен превышать 20 МБ"})
		return
	}

	// Проверка типа файла
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл должен быть изображением"})
		return
	}

	// Генерируем уникальное имя файла
	uniqueFileName := generateUniqueFileName(file.Filename)
	imagePath := fmt.Sprintf("profile-backgrounds/%s/%s", user.Nickname, uniqueFileName)

	// Создаем временную директорию
	tempDir := "./temp_uploads"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания временной директории"})
		return
	}

	// Сохраняем файл временно
	tempFilePath := filepath.Join(tempDir, uniqueFileName)
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения файла"})
		return
	}
	defer os.Remove(tempFilePath)

	// Удаляем старое изображение из S3, если оно есть
	if user.BackgroundAvatarProfile != "" {
		oldKey := extractS3KeyFromURL(user.BackgroundAvatarProfile)
		if oldKey != "" {
			if err := DeleteFileFromS3(oldKey); err != nil {
				fmt.Printf("Предупреждение: не удалось удалить старый файл из S3: %v\n", err)
				// Продолжаем выполнение, даже если удаление не удалось
			}
		}
	}

	// Загружаем новое изображение в S3
	fmt.Println("Загружаем фон профиля в S3")
	publicURL, errS3 := UploadFileToS3(tempFilePath, imagePath)
	if errS3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка загрузки в S3: %v", errS3)})
		return
	}
	fmt.Println("Фон профиля загружен:", publicURL)

	// Обновляем БД через сервисный слой
	if err := services.UpdateProfileBackground(nickname.(string), publicURL); err != nil {
		// Если обновление БД не удалось, удаляем только что загруженный файл
		DeleteFileFromS3(imagePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления профиля"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":                   "Фон профиля успешно обновлен",
		"background_avatar_profile": publicURL,
	})
}

// Удаление фона профиля
func DeleteProfileBackground(c *gin.Context) {
	// Получаем nickname из контекста
	nickname, exists := c.Get("nickname")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	// Получаем текущего пользователя
	user, err := services.GetUserByNickname(nickname.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных пользователя"})
		return
	}

	// Проверяем, есть ли фон для удаления
	if user.BackgroundAvatarProfile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Фон профиля уже отсутствует"})
		return
	}

	// Удаляем файл из S3
	oldKey := extractS3KeyFromURL(user.BackgroundAvatarProfile)
	if oldKey != "" {
		if err := DeleteFileFromS3(oldKey); err != nil {
			fmt.Printf("Предупреждение: не удалось удалить файл из S3: %v\n", err)
			// Продолжаем, даже если удаление из S3 не удалось
		}
	}

	// Обновляем БД через сервисный слой - устанавливаем пустую строку
	if err := services.DeleteProfileBackground(nickname.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления профиля"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Фон профиля успешно удален",
	})
}

// UpdateTelegram обновляет telegram пользователя
func UpdateTelegram(c *gin.Context) {
	nickname, exists := c.Get("nickname")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	var req struct {
		Telegram string `json:"telegram" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Валидация telegram (должен начинаться с @ или быть username)
	telegram := strings.TrimSpace(req.Telegram)
	if telegram == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Telegram не может быть пустым"})
		return
	}

	// Добавляем @ если отсутствует
	if !strings.HasPrefix(telegram, "@") && !strings.HasPrefix(telegram, "https://t.me/") {
		telegram = "@" + telegram
	}

	// Обновляем telegram через сервисный слой
	if err := services.UpdateTelegram(nickname.(string), telegram); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления telegram"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Telegram успешно обновлен",
		"telegram": telegram,
	})
}
