package models

import "time"

type Ad struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ServerName       string    `gorm:"column:server_name" json:"server_name"`
	Title            string    `gorm:"column:title" json:"title"`
	Description      string    `gorm:"column:description" json:"description"`
	Type             string    `gorm:"column:type" json:"type"`
	Currency         *string   `gorm:"column:currency" json:"currency,omitempty"`
	Price            *int64    `gorm:"column:price" json:"price,omitempty"`
	PricePeriod      *string   `gorm:"column:price_period" json:"price_period,omitempty"`
	RentalHoursLimit *int      `gorm:"column:rental_hours_limit" json:"rental_hours_limit,omitempty"`
	Image            string    `gorm:"column:image" json:"image"`
	Category         string    `gorm:"column:category" json:"category"`
	Nickname         string    `gorm:"column:nickname" json:"nickname"`
	Views            int       `gorm:"column:views;default:0" json:"views"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

type Report struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AdID             uint      `gorm:"column:ad_id;not null" json:"ad_id"`
	ReporterNickname string    `gorm:"column:reporter_nickname;size:50;not null" json:"reporter_nickname"`
	Reason           string    `gorm:"column:reason;size:100;not null" json:"reason"`
	Description      *string   `gorm:"column:description;type:text" json:"description,omitempty"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}
