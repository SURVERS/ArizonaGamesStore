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

// CreateNewAds godoc
// @Summary Создать объявление
// @Description Создает новое объявление с картинкой. Картинка загружается на AWS S3. После создания объявление автоматически удалится через 48 часов (можно продлить). Между созданиями объявлений нужно ждать 60 секунд
// @Tags Объявления
// @Accept multipart/form-data
// @Produce json
// @Param server formData string true "Сервер (ViceCity, Phoenix, и т.д.)"
// @Param title formData string true "Название (макс. 25 символов)"
// @Param description formData string true "Описание (макс. 500 символов)"
// @Param type formData string true "Тип (Продать/Купить/Сдать в аренду)"
// @Param currency formData string true "Валюта (VC/$/BTC/EURO/Договорная)"
// @Param price formData number true "Цена"
// @Param category formData string true "Категория (house/business/vehicle/security/accs/others)"
// @Param nickname formData string true "Никнейм автора"
// @Param imagePath formData string true "Путь для сохранения картинки"
// @Param image formData file true "Изображение (макс. 10MB, разрешение 300x200 - 1920x1080)"
// @Param rentalHoursLimit formData int false "Лимит часов аренды (1-180)"
// @Success 200 {object} map[string]interface{} "Объявление создано! ID: 42"
// @Failure 400 {object} map[string]string "Не хватает данных или картинка кривая"
// @Failure 413 {object} map[string]string "Картинка слишком большая (макс. 10MB)"
// @Failure 429 {object} map[string]string "Подожди 60 секунд перед созданием нового объявления"
// @Failure 500 {object} map[string]string "Ошибка загрузки на S3 или БД"
// @Router /createnewads [post]
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

// GetAdsByCategory godoc
// @Summary Список объявлений
// @Description Возвращает список объявлений с фильтрацией и сортировкой. По умолчанию возвращает 20 штук, можно подгружать дальше через offset
// @Tags Объявления
// @Produce json
// @Param category query string true "Категория" Enums(house, business, vehicle, security, accs, others)
// @Param server query string false "Фильтр по серверу"
// @Param limit query int false "Сколько объявлений вернуть (по умолчанию 20)"
// @Param offset query int false "Сколько пропустить для пагинации (по умолчанию 0)"
// @Param sort query string false "Сортировка" Enums(date_desc, date_asc, price_desc, price_asc, views_desc)
// @Param type query string false "Фильтр по типу" Enums(Продать, Купить, Сдать в аренду)
// @Param currency query string false "Фильтр по валюте" Enums(VC, $, BTC, EURO, Договорная)
// @Param price_min query number false "Минимальная цена"
// @Param price_max query number false "Максимальная цена"
// @Success 200 {object} map[string]interface{} "Список объявлений"
// @Failure 400 {object} map[string]string "Не указана категория"
// @Failure 500 {object} map[string]string "Ошибка БД"
// @Router /ads [get]
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

// GetAdsByNickname godoc
// @Summary Объявления по нику
// @Description Возвращает все объявления конкретного пользователя
// @Tags Объявления
// @Produce json
// @Param nickname path string true "Никнейм пользователя"
// @Success 200 {object} map[string]interface{} "Список объявлений пользователя"
// @Failure 500 {object} map[string]string "Ошибка БД"
// @Router /listings/user/{nickname} [get]
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

// CreateReport godoc
// @Summary Пожаловаться на объявление
// @Description Отправляет жалобу на объявление. Причины: Мошенничество, Спам, Порнография, и т.д. Жалобы проверяются модераторами
// @Tags Жалобы
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Данные жалобы"
// @Success 200 {object} map[string]string "Жалоба отправлена! Мы её рассмотрим"
// @Failure 400 {object} map[string]string "Не указана причина или описание слишком короткое"
// @Failure 401 {object} map[string]string "Залогинься чтобы жаловаться"
// @Failure 404 {object} map[string]string "Объявление не найдено"
// @Failure 500 {object} map[string]string "Ошибка создания жалобы"
// @Router /reports [post]
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

// GetRandomAds godoc
// @Summary Случайные объявления
// @Description Возвращает 8 случайных объявлений для главной страницы. Каждый раз разные!
// @Tags Объявления
// @Produce json
// @Success 200 {object} map[string]interface{} "Массив случайных объявлений"
// @Failure 500 {object} map[string]string "Что-то пошло не так"
// @Router /ads/random [get]
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

// UpdateAd godoc
// @Summary Обновить объявление
// @Description Обновляет данные объявления. Доступно только автору объявления
// @Tags Объявления
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID объявления"
// @Param request body map[string]interface{} true "Новые данные"
// @Success 200 {object} map[string]interface{} "Объявление обновлено!"
// @Failure 400 {object} map[string]string "Некорректные данные"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 403 {object} map[string]string "Это не твое объявление!"
// @Failure 404 {object} map[string]string "Объявление не найдено"
// @Failure 500 {object} map[string]string "Ошибка обновления"
// @Router /ads/{id} [put]
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

// DeleteAd godoc
// @Summary Удалить объявление
// @Description Удаляет объявление. Доступно только автору
// @Tags Объявления
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID объявления"
// @Success 200 {object} map[string]string "Объявление удалено!"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 403 {object} map[string]string "Это не твое объявление!"
// @Failure 404 {object} map[string]string "Объявление не найдено"
// @Failure 500 {object} map[string]string "Ошибка удаления"
// @Router /ads/{id} [delete]
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

// IncrementAdViews godoc
// @Summary Записать просмотр
// @Description Увеличивает счетчик просмотров объявления. Вызывай когда пользователь открывает карточку объявления
// @Tags Объявления
// @Produce json
// @Param id path int true "ID объявления"
// @Success 200 {object} map[string]string "Просмотр засчитан!"
// @Failure 404 {object} map[string]string "Объявление не найдено"
// @Failure 500 {object} map[string]string "Ошибка обновления"
// @Router /ads/{id}/view [post]
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
