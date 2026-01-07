package models

import (
	"time"
)

type FeedbackAd struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AdID             int       `gorm:"not null" json:"ad_id"`
	ReviewerNickname string    `gorm:"not null" json:"reviewer_nickname"`
	AdOwnerNickname  string    `gorm:"not null" json:"ad_owner_nickname"`
	Rating           int       `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	ReviewText       string    `gorm:"type:text;not null" json:"review_text"`
	ProofImage       string    `gorm:"not null" json:"proof_image"`
	ConfirmFeedback  bool      `gorm:"default:false" json:"confirm_feedback"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (FeedbackAd) TableName() string {
	return "feedback_ads"
}

type ViewedAd struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserNickname string    `gorm:"not null" json:"user_nickname"`
	AdID         int       `gorm:"not null" json:"ad_id"`
	ViewedAt     time.Time `gorm:"autoCreateTime" json:"viewed_at"`
}

func (ViewedAd) TableName() string {
	return "viewed_ads"
}

type FeedbackWithReviewer struct {
	ID               uint      `json:"id"`
	AdID             int       `json:"ad_id"`
	ReviewerNickname string    `json:"reviewer_nickname"`
	ReviewerAvatar   string    `json:"reviewer_avatar"`
	ReviewerRating   float32   `json:"reviewer_rating"`
	AdOwnerNickname  string    `json:"ad_owner_nickname"`
	Rating           int       `json:"rating"`
	ReviewText       string    `json:"review_text"`
	ProofImage       string    `json:"proof_image"`
	ConfirmFeedback  bool      `json:"confirm_feedback"`
	CreatedAt        time.Time `json:"created_at"`
}
