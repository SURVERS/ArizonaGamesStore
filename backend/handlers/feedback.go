package handlers

import (
	"arizonagamesstore/backend/database"
	"arizonagamesstore/backend/models"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateFeedback godoc
// @Summary Оставить отзыв
// @Description Оставляет отзыв о продавце. Рейтинг от 1 до 5 звезд. Отзыв нужно подтвердить продавцу, чтобы он отобразился
// @Tags Отзывы
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Данные отзыва"
// @Success 200 {object} map[string]interface{} "Отзыв отправлен! Ждем подтверждения продавца"
// @Failure 400 {object} map[string]string "Некорректный рейтинг или комментарий"
// @Failure 401 {object} map[string]string "Нужна авторизация"
// @Failure 409 {object} map[string]string "Ты уже оставлял отзыв этому продавцу"
// @Failure 500 {object} map[string]string "Ошибка сохранения"
// @Router /feedback [post]
func CreateFeedback(c *gin.Context) {
	reviewerNickname, exists := c.Get("nickname")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	// Получение файла изображения
	file, err := c.FormFile("proof_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Изображение-доказательство обязательно"})
		return
	}

	// Проверка размера файла (15MB)
	if file.Size > 15*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Размер изображения не должен превышать 15 МБ"})
		return
	}

	adID, err := strconv.Atoi(c.PostForm("ad_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID объявления"})
		return
	}

	rating, err := strconv.Atoi(c.PostForm("rating"))
	if err != nil || rating < 1 || rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Рейтинг должен быть от 1 до 5"})
		return
	}

	reviewText := c.PostForm("review_text")
	if reviewText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Текст отзыва обязателен"})
		return
	}

	// Получить информацию о объявлении
	var ad models.Ad
	if err := database.DB.Where("id = ?", adID).First(&ad).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Объявление не найдено"})
		return
	}

	// Проверить, что пользователь не оставляет отзыв на свое объявление
	if ad.Nickname == reviewerNickname.(string) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нельзя оставить отзыв на свое объявление"})
		return
	}

	// Проверить, не оставлял ли пользователь уже отзыв на это объявление
	var existingFeedback models.FeedbackAd
	result := database.DB.Where("ad_id = ? AND reviewer_nickname = ?", adID, reviewerNickname.(string)).First(&existingFeedback)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Вы уже оставили отзыв на это объявление"})
		return
	}

	// Сохранение изображения
	uploadsDir := "./uploads/feedbacks"
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания директории"})
		return
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", reviewerNickname.(string), time.Now().Unix(), ext)
	imagePath := filepath.Join(uploadsDir, filename)

	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения изображения"})
		return
	}

	// Сохраняем URL для доступа через статическую папку /uploads
	imageURL := fmt.Sprintf("http://localhost:8080/uploads/feedbacks/%s", filename)

	// Создание отзыва
	feedback := models.FeedbackAd{
		AdID:             adID,
		ReviewerNickname: reviewerNickname.(string),
		AdOwnerNickname:  ad.Nickname,
		Rating:           rating,
		ReviewText:       reviewText,
		ProofImage:       imageURL,
		ConfirmFeedback:  false,
	}

	if err := database.DB.Create(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка создания отзыва: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Отзыв успешно отправлен на модерацию",
		"feedback": feedback,
	})
}

// GetFeedbacksByOwner godoc
// @Summary Отзывы продавца
// @Description Возвращает все отзывы о продавце
// @Tags Отзывы
// @Produce json
// @Param nickname path string true "Никнейм продавца"
// @Success 200 {object} map[string]interface{} "Список отзывов"
// @Failure 500 {object} map[string]string "Ошибка загрузки"
// @Router /feedback/{nickname} [get]
func GetFeedbacksByOwner(c *gin.Context) {
	ownerNickname := c.Param("nickname")

	var feedbacks []models.FeedbackWithReviewer

	result := database.DB.Table("feedback_ads").
		Select("feedback_ads.*, accounts.avatar as reviewer_avatar, accounts.rating as reviewer_rating").
		Joins("LEFT JOIN accounts ON feedback_ads.reviewer_nickname = accounts.nickname").
		Where("feedback_ads.ad_owner_nickname = ? AND feedback_ads.confirm_feedback = ?", ownerNickname, true).
		Order("feedback_ads.created_at DESC").
		Find(&feedbacks)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка получения отзывов: %v", result.Error)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"feedbacks": feedbacks})
}

// ConfirmFeedback godoc
// @Summary Подтвердить отзыв
// @Description Подтверждает отзыв. Только продавец может подтвердить отзыв о себе. После подтверждения рейтинг пересчитывается
// @Tags Отзывы
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID отзыва"
// @Success 200 {object} map[string]string "Отзыв подтвержден! Твой рейтинг обновлен"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 403 {object} map[string]string "Ты не можешь подтвердить этот отзыв"
// @Failure 404 {object} map[string]string "Отзыв не найден"
// @Failure 500 {object} map[string]string "Ошибка подтверждения"
// @Router /feedback/{id}/confirm [put]
func ConfirmFeedback(c *gin.Context) {
	feedbackID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID отзыва"})
		return
	}

	var feedback models.FeedbackAd
	if err := database.DB.Where("id = ?", feedbackID).First(&feedback).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Отзыв не найден"})
		return
	}

	feedback.ConfirmFeedback = true
	if err := database.DB.Save(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подтверждения отзыва"})
		return
	}

	// Пересчитываем рейтинг владельца объявления
	if err := UpdateUserRating(feedback.AdOwnerNickname); err != nil {
		// Логируем ошибку, но не возвращаем её клиенту, так как отзыв уже подтвержден
		fmt.Printf("Ошибка обновления рейтинга пользователя %s: %v\n", feedback.AdOwnerNickname, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Отзыв подтвержден"})
}

// AddViewedAd godoc
// @Summary Добавить в историю
// @Description Добавляет объявление в историю просмотров пользователя
// @Tags Просмотренное
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]int true "ID объявления" example(ad_id=42)
// @Success 200 {object} map[string]string "Добавлено в историю"
// @Failure 400 {object} map[string]string "Не указан ID"
// @Failure 401 {object} map[string]string "Нужна авторизация"
// @Failure 500 {object} map[string]string "Ошибка сохранения"
// @Router /viewed-ads [post]
func AddViewedAd(c *gin.Context) {
	userNickname, exists := c.Get("nickname")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	adID, err := strconv.Atoi(c.PostForm("ad_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID объявления"})
		return
	}

	// Проверить существование объявления
	var ad models.Ad
	if err := database.DB.Where("id = ?", adID).First(&ad).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Объявление не найдено"})
		return
	}

	// Не добавлять свои объявления в просмотренные
	if ad.Nickname == userNickname.(string) {
		c.JSON(http.StatusOK, gin.H{"message": "Это ваше объявление"})
		return
	}

	// Проверить, не просмотрено ли уже
	var existingView models.ViewedAd
	result := database.DB.Where("user_nickname = ? AND ad_id = ?", userNickname.(string), adID).First(&existingView)

	if result.Error == nil {
		// Уже просмотрено, обновляем время
		existingView.ViewedAt = time.Now()
		database.DB.Save(&existingView)
	} else {
		// Создаем новую запись
		viewedAd := models.ViewedAd{
			UserNickname: userNickname.(string),
			AdID:         adID,
		}
		if err := database.DB.Create(&viewedAd).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения просмотра"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Объявление добавлено в просмотренные"})
}

// GetViewedAds godoc
// @Summary История просмотров
// @Description Возвращает историю просмотренных объявлений пользователя
// @Tags Просмотренное
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Список просмотренных объявлений"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 500 {object} map[string]string "Ошибка загрузки"
// @Router /viewed-ads [get]
func GetViewedAds(c *gin.Context) {
	userNickname, exists := c.Get("nickname")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	// Получаем viewed_ads с соответствующими объявлениями
	var viewedRecords []models.ViewedAd
	result := database.DB.Where("user_nickname = ?", userNickname.(string)).
		Order("viewed_at DESC").
		Find(&viewedRecords)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка получения просмотренных: %v", result.Error)})
		return
	}

	// Создаем ответ с полной информацией об объявлениях
	type AdWithAuthor struct {
		models.Ad
		AuthorAvatar  string  `json:"author_avatar"`
		AuthorRating  float32 `json:"author_rating"`
		OwnerTelegram string  `json:"owner_telegram"`
	}

	type ViewedAdResponse struct {
		ID           uint          `json:"id"`
		UserNickname string        `json:"user_nickname"`
		AdID         int           `json:"ad_id"`
		ViewedAt     time.Time     `json:"viewed_at"`
		Ad           *AdWithAuthor `json:"Ad"`
	}

	viewedAds := make([]ViewedAdResponse, 0, len(viewedRecords))
	for _, viewed := range viewedRecords {
		var ad AdWithAuthor
		err := database.DB.Table("ads").
			Select("ads.*, accounts.avatar as author_avatar, accounts.rating as author_rating, accounts.telegram as owner_telegram").
			Joins("LEFT JOIN accounts ON ads.nickname = accounts.nickname").
			Where("ads.id = ?", viewed.AdID).
			First(&ad).Error

		if err == nil {
			viewedAds = append(viewedAds, ViewedAdResponse{
				ID:           viewed.ID,
				UserNickname: viewed.UserNickname,
				AdID:         viewed.AdID,
				ViewedAt:     viewed.ViewedAt,
				Ad:           &ad,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"viewed_ads": viewedAds})
}

// UpdateUserRating пересчитывает средний рейтинг пользователя на основе подтвержденных отзывов
func UpdateUserRating(nickname string) error {
	// Получаем все подтвержденные отзывы для данного пользователя
	var feedbacks []models.FeedbackAd
	result := database.DB.Where("ad_owner_nickname = ? AND confirm_feedback = ?", nickname, true).Find(&feedbacks)

	if result.Error != nil {
		return result.Error
	}

	// Если нет отзывов, устанавливаем рейтинг 0
	if len(feedbacks) == 0 {
		return database.DB.Model(&models.Account{}).
			Where("nickname = ?", nickname).
			Update("rating", 0).Error
	}

	// Вычисляем средний рейтинг
	var totalRating int
	for _, feedback := range feedbacks {
		totalRating += feedback.Rating
	}
	averageRating := float32(totalRating) / float32(len(feedbacks))

	// Обновляем рейтинг пользователя
	return database.DB.Model(&models.Account{}).
		Where("nickname = ?", nickname).
		Update("rating", averageRating).Error
}
