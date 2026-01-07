package addcells

import (
	"arizonagamesstore/backend/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedStatistics(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Statistic{}).Count(&count).Error; err != nil {
		if err := db.AutoMigrate(&models.Statistic{}); err != nil {
			return err
		}
	}

	if count > 0 {
		return nil
	}

	initialStats := []models.Statistic{
		{CategoryName: "accs", AdCount: 0},
		{CategoryName: "business", AdCount: 0},
		{CategoryName: "house", AdCount: 0},
		{CategoryName: "security", AdCount: 0},
		{CategoryName: "vehicles", AdCount: 0},
		{CategoryName: "others", AdCount: 0},
	}

	for _, stat := range initialStats {
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "category_name"}},
			DoNothing: true,
		}).Create(&stat).Error; err != nil {
			return err
		}
	}

	return nil
}
