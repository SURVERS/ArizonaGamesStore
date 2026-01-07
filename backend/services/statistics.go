package services

import (
	"arizonagamesstore/backend/database"
	"arizonagamesstore/backend/models"

	"gorm.io/gorm"
)

func AddAdCount(category_name string) error {
	err := database.DB.Model(&models.Statistic{}).
		Where("category_name = ?", category_name).
		UpdateColumn("ad_count", gorm.Expr("ad_count + ?", 1))
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func GetAdCounts(category_name string) (int, error) {
	var adCount int
	err := database.DB.Model(&models.Statistic{}).
		Select("ad_count").
		Where("category_name = ?", category_name).
		Scan(&adCount).Error
	if err != nil {
		return 0, err
	}
	return adCount, nil
}

func DecreaseAdCount(category_name string) error {
	err := database.DB.Model(&models.Statistic{}).
		Where("category_name = ?", category_name).
		UpdateColumn("ad_count", gorm.Expr("GREATEST(ad_count - 1, 0)"))
	if err.Error != nil {
		return err.Error
	}
	return nil
}
