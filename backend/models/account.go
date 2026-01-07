package models

import (
	"time"
)

type Account struct {
	ID                      uint       `gorm:"primaryKey;autoIncrement"`
	Nickname                string     `gorm:"uniqueIndex;not null"`
	Email                   string     `gorm:"uniqueIndex"`
	EmailVerified           bool       `gorm:"column:email_verified;default:false"`
	PasswordHash            string     `gorm:"column:password_hash"`
	UserRole                string     `gorm:"column:user_role"`
	Avatar                  string     `gorm:"column:avatar"`
	BackgroundAvatarProfile string     `gorm:"column:background_avatar_profile"`
	Rating                  float32    `gorm:"column:rating"`
	SuccessTransactions     int        `gorm:"column:success_transactions"`
	UserDescription         string     `gorm:"column:user_description"`
	Theme                   string     `gorm:"column:theme;default:'dark'"`
	Telegram                string     `gorm:"column:telegram"`
	LastSeenAt              time.Time  `gorm:"column:last_seen_at"`
	LastSettingsChange      *time.Time `gorm:"column:last_settings_change"`
	LastNicknameChange      *time.Time `gorm:"column:last_nickname_change"`
	LastEmailChange         *time.Time `gorm:"column:last_email_change"`
	RegIP                   string     `gorm:"column:reg_ip"`
	LastIP                  string     `gorm:"column:last_ip"`
	CreatedAt               time.Time  `gorm:"autoCreateTime"`
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	AccountID uint      `gorm:"not null;index"`
	Token     string    `gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type EmailVerification struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Email        string    `gorm:"not null"`
	Nickname     string    `gorm:"not null"`
	PasswordHash string    `gorm:"not null"`
	Code         string    `gorm:"not null"`
	ExpiresAt    time.Time `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
