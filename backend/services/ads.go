package services

import (
	"arizonagamesstore/backend/database"
	"arizonagamesstore/backend/models"
	"fmt"
	"log"
	"time"
)

func CreateNewAd(dto models.Ad, filePathS3 string) (bool, string) {
	createAd := models.Ad{
		ServerName:       dto.ServerName,
		Title:            dto.Title,
		Description:      dto.Description,
		Type:             dto.Type,
		Currency:         dto.Currency,
		Price:            dto.Price,
		PricePeriod:      dto.PricePeriod,
		RentalHoursLimit: dto.RentalHoursLimit,
		Category:         dto.Category,
		Nickname:         dto.Nickname,
		Image:            filePathS3,
	}

	result := database.DB.Create(&createAd)
	if result.Error != nil {
		err := fmt.Sprint(result.Error)
		return false, err
	}

	return true, ""
}

type AdWithAuthor struct {
	models.Ad
	AuthorAvatar   string  `json:"author_avatar"`
	AuthorRating   float32 `json:"author_rating"`
	OwnerTelegram  string  `json:"owner_telegram"`
}

type AdFilters struct {
	Sort     string
	Type     string
	PriceMin *float64
	PriceMax *float64
	Currency string
}

func GetAdsByCategory(category string, server string, limit int, offset int, filters *AdFilters) ([]AdWithAuthor, error) {
	var ads []AdWithAuthor

	query := database.DB.Table("ads").
		Select("ads.*, accounts.avatar as author_avatar, accounts.rating as author_rating, accounts.telegram as owner_telegram").
		Joins("LEFT JOIN accounts ON ads.nickname = accounts.nickname").
		Where("ads.category = ?", category)

	if server != "" && server != "all" {
		query = query.Where("ads.server_name = ?", server)
	}

	if filters != nil {
		if filters.Type != "" {
			query = query.Where("ads.type = ?", filters.Type)
		}
		if filters.PriceMin != nil {
			query = query.Where("ads.price >= ?", *filters.PriceMin)
		}
		if filters.PriceMax != nil {
			query = query.Where("ads.price <= ?", *filters.PriceMax)
		}
		if filters.Currency != "" {
			query = query.Where("ads.currency = ?", filters.Currency)
		}

		switch filters.Sort {
		case "date_asc":
			query = query.Order("ads.created_at ASC")
		case "price_desc":
			query = query.Order("ads.price DESC")
		case "price_asc":
			query = query.Order("ads.price ASC")
		case "views_desc":
			query = query.Order("ads.views DESC")
		default:
			query = query.Order("ads.created_at DESC")
		}
	} else {
		query = query.Order("ads.created_at DESC")
	}

	result := query.Limit(limit).Offset(offset).Find(&ads)
	if result.Error != nil {
		return nil, result.Error
	}

	return ads, nil
}

func GetAdsByNickname(nickname string) ([]AdWithAuthor, error) {
	var ads []AdWithAuthor

	result := database.DB.Table("ads").
		Select("ads.*, accounts.avatar as author_avatar, accounts.rating as author_rating, accounts.telegram as owner_telegram").
		Joins("LEFT JOIN accounts ON ads.nickname = accounts.nickname").
		Where("ads.nickname = ?", nickname).
		Order("ads.created_at DESC").
		Find(&ads)

	if result.Error != nil {
		return nil, result.Error
	}

	return ads, nil
}

func UpdateNickNameAds(oldNickname string, newNickname string) error {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&models.Ad{}).
		Where("nickname = ?", oldNickname).
		Updates(map[string]interface{}{
			"nickname": newNickname,
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func CreateReport(adID uint, reporterNickname string, reason string, description *string) error {
	report := models.Report{
		AdID:             adID,
		ReporterNickname: reporterNickname,
		Reason:           reason,
		Description:      description,
	}

	result := database.DB.Create(&report)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func AutoDeleteOldAds() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	log.Println("Запущена служба автоудаления объявлений старше 48 часов")

	for {
		<-ticker.C

		cutoffTime := time.Now().Add(-48 * time.Hour)

		var adsToDelete []models.Ad
		database.DB.Where("created_at < ?", cutoffTime).Find(&adsToDelete)

		categoryCount := make(map[string]int)
		for _, ad := range adsToDelete {
			categoryCount[ad.Category]++
		}

		result := database.DB.Where("created_at < ?", cutoffTime).Delete(&models.Ad{})

		if result.Error != nil {
			log.Printf("Ошибка при удалении старых объявлений: %v", result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("Удалено %d объявлений старше 48 часов", result.RowsAffected)

			for category, count := range categoryCount {
				for i := 0; i < count; i++ {
					DecreaseAdCount(category)
				}
				log.Printf("Уменьшен счетчик для категории %s на %d", category, count)
			}
		}
	}
}

func RecalculateStatistics() error {
	var stats []models.Statistic
	if err := database.DB.Find(&stats).Error; err != nil {
		return err
	}

	for _, stat := range stats {
		var count int64
		database.DB.Table("ads").Where("category = ?", stat.CategoryName).Count(&count)

		database.DB.Model(&models.Statistic{}).
			Where("category_name = ?", stat.CategoryName).
			Update("ad_count", count)

		log.Printf("Пересчет счетчика для %s: %d объявлений", stat.CategoryName, count)
	}

	return nil
}

func GetRandomAds(limit int, offset int) ([]AdWithAuthor, error) {
	var ads []AdWithAuthor

	result := database.DB.Table("ads").
		Select("ads.*, accounts.avatar as author_avatar, accounts.rating as author_rating, accounts.telegram as owner_telegram").
		Joins("LEFT JOIN accounts ON ads.nickname = accounts.nickname").
		Order("RANDOM()").
		Limit(limit).
		Offset(offset).
		Find(&ads)

	if result.Error != nil {
		return nil, result.Error
	}

	return ads, nil
}
