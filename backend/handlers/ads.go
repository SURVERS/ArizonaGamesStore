package handlers

import (
	"arizonagamesstore/backend/database"
	"arizonagamesstore/backend/models"
	"arizonagamesstore/backend/services"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type NewAddsRequest struct {
	Server           string  `form:"server" binding:"required"`
	Title            string  `form:"title" binding:"required"`
	Description      string  `form:"description" binding:"required"`
	Types            string  `form:"type" binding:"required"`
	Currency         *string `form:"currency"`
	Price            *int64  `form:"price"`
	PricePeriod      *string `form:"pricePeriod"`
	RentalHoursLimit *int    `form:"rentalHoursLimit"`
	Category         string  `form:"category" binding:"required"`
	Nickname         string  `form:"nickname" binding:"required"`
	ImagePath        string  `form:"imagePath" binding:"required"`
}

func CreateNewAds(c *gin.Context) {
	var req NewAddsRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получение файла изображения
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Изображение обязательно"})
		return
	}

	// Проверка размера файла (10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Размер изображения не должен превышать 10 МБ"})
		return
	}

	dto := models.Ad{
		ServerName:       req.Server,
		Title:            req.Title,
		Description:      req.Description,
		Type:             req.Types,
		Currency:         req.Currency,
		Price:            req.Price,
		PricePeriod:      req.PricePeriod,
		RentalHoursLimit: req.RentalHoursLimit,
		Category:         req.Category,
		Nickname:         req.Nickname,
	}

	fmt.Println("Server:", req.Server)
	fmt.Println("Title:", req.Title)
	fmt.Println("Description:", req.Description)
	fmt.Println("Type:", req.Types)
	fmt.Println("Category:", req.Category)
	fmt.Println("Nickname:", req.Nickname)
	fmt.Println("Image filename:", file.Filename)
	fmt.Println("Image path:", req.ImagePath)
	fmt.Println("Image size:", file.Size)

	// Создаем временную директорию если её нет
	tempDir := "./temp_uploads"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания временной директории"})
		return
	}

	// Сохраняем файл во временную директорию
	tempFilePath := filepath.Join(tempDir, file.Filename)
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения файла"})
		return
	}
	defer os.Remove(tempFilePath) // Удаляем временный файл после загрузки

	fmt.Println("Отправляем в S3 хранилище изображение")
	publicURL, errS3 := UploadFileToS3(tempFilePath, req.ImagePath)
	if errS3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка при загрузке изображения: %v", errS3)})
		return
	}
	fmt.Println("Изображение успешно загружено в S3")
	result, errDB := services.CreateNewAd(dto, publicURL)
	if result != true {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка при создании объявления: %v", errDB)})
		return
	}

	errUpdate := services.AddAdCount(dto.Category)
	if errUpdate != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка при обновлении статистики: %v", errUpdate)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Объявление успешно создано"})
}

// GetAdsByCategory обработчик для получения объявлений по категории
func GetAdsByCategory(c *gin.Context) {
	category := c.Query("category")
	server := c.Query("server")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Категория обязательна"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Парсим параметры фильтрации и сортировки
	filters := &services.AdFilters{
		Sort:     c.Query("sort"),
		Type:     c.Query("type"),
		Currency: c.Query("currency"),
	}

	// Парсим диапазон цен
	if priceMinStr := c.Query("price_min"); priceMinStr != "" {
		if priceMin, err := strconv.ParseFloat(priceMinStr, 64); err == nil {
			filters.PriceMin = &priceMin
		}
	}
	if priceMaxStr := c.Query("price_max"); priceMaxStr != "" {
		if priceMax, err := strconv.ParseFloat(priceMaxStr, 64); err == nil {
			filters.PriceMax = &priceMax
		}
	}

	ads, err := services.GetAdsByCategory(category, server, limit, offset, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка получения объявлений: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ads": ads})
}

// GetAdsByNickname обработчик для получения объявлений пользователя
func GetAdsByNickname(c *gin.Context) {
	nickname := c.Param("nickname")

	if nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Никнейм обязателен"})
		return
	}

	ads, err := services.GetAdsByNickname(nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка получения объявлений: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"listings": ads,
	})
}

// CreateReport обработчик для создания жалобы на объявление
func CreateReport(c *gin.Context) {
	var req struct {
		AdID        uint    `json:"ad_id" binding:"required"`
		Reason      string  `json:"reason" binding:"required"`
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	// Получаем nickname из контекста (если используется middleware)
	nickname, exists := c.Get("nickname")
	if !exists {
		// Если нет middleware, можно использовать из body или другой источник
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Необходима авторизация"})
		return
	}

	err := services.CreateReport(req.AdID, nickname.(string), req.Reason, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка создания жалобы: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Жалоба успешно отправлена"})
}

// GetRandomAds обработчик для получения рандомных объявлений из разных категорий
func GetRandomAds(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "15")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 15
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	ads, err := services.GetRandomAds(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка получения объявлений: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ads": ads})
}

// UpdateAd обновляет объявление
func UpdateAd(c *gin.Context) {
	nickname, exists := c.Get("nickname")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	adID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID объявления"})
		return
	}

	// Получить объявление и проверить владельца
	var ad models.Ad
	if err := database.DB.Where("id = ?", adID).First(&ad).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Объявление не найдено"})
		return
	}

	if ad.Nickname != nickname.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете редактировать это объявление"})
		return
	}

	// Обновление полей
	if title := c.PostForm("title"); title != "" {
		ad.Title = title
	}
	if description := c.PostForm("description"); description != "" {
		ad.Description = description
	}
	if typeStr := c.PostForm("type"); typeStr != "" {
		ad.Type = typeStr
	}
	if priceStr := c.PostForm("price"); priceStr != "" {
		price, err := strconv.ParseInt(priceStr, 10, 64)
		if err == nil {
			ad.Price = &price
		}
	}

	// Обработка изображения, если оно предоставлено
	file, err := c.FormFile("image")
	if err == nil {
		// Проверка размера файла (10MB)
		if file.Size > 10*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Размер изображения не должен превышать 10 МБ"})
			return
		}

		// Удаление старого изображения
		if ad.Image != "" {
			os.Remove(ad.Image)
		}

		// Сохранение нового изображения
		uploadsDir := "./uploads/ads"
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%s_%d%s", nickname.(string), time.Now().Unix(), ext)
		imagePath := filepath.Join(uploadsDir, filename)

		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения изображения"})
			return
		}

		ad.Image = imagePath
	}

	// Сохранение изменений
	if err := database.DB.Save(&ad).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления объявления"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Объявление успешно обновлено",
		"ad":      ad,
	})
}

// DeleteAd удаляет объявление
func DeleteAd(c *gin.Context) {
	nickname, exists := c.Get("nickname")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	adID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID объявления"})
		return
	}

	// Получить объявление и проверить владельца
	var ad models.Ad
	if err := database.DB.Where("id = ?", adID).First(&ad).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Объявление не найдено"})
		return
	}

	if ad.Nickname != nickname.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете удалить это объявление"})
		return
	}

	// Удаление изображения
	if ad.Image != "" {
		os.Remove(ad.Image)
	}

	// Удаление объявления
	if err := database.DB.Delete(&ad).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления объявления"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Объявление успешно удалено"})
}

// IncrementAdViews увеличивает счетчик просмотров объявления
func IncrementAdViews(c *gin.Context) {
	adID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID объявления"})
		return
	}

	// Увеличение счетчика просмотров
	if err := database.DB.Model(&models.Ad{}).Where("id = ?", adID).UpdateColumn("views", database.DB.Raw("views + 1")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления просмотров"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Просмотр учтен"})
}
