package main

import (
	_ "arizonagamesstore/backend/docs"
	"arizonagamesstore/backend/database"
	"arizonagamesstore/backend/handlers"
	"arizonagamesstore/backend/middleware"
	"arizonagamesstore/backend/services"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"log"

	"github.com/joho/godotenv"
)

// @title Arizona Games Store API
// @version 1.0
// @description API для игрового маркетплейса Arizona RP. Здесь можно купить/продать/арендовать дома, бизнесы, транспорт и всякую другую всячину. Работает на 33 серверах, поддерживает несколько валют и умеет в рейтинги продавцов.
// @description
// @description Основные фишки:
// @description - JWT авторизация (access + refresh токены)
// @description - Rate limiting чтобы боты не спамили
// @description - Email верификация
// @description - Загрузка картинок в AWS S3
// @description - Система отзывов и рейтингов
// @description - Жалобы на объявления
// @description - Автоудаление старых объявлений через 48 часов
// @description
// @description Сделано с душой и большим количеством кофе ☕

// @contact.name Поддержка
// @contact.email support@arizonagamesstore.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Вставь сюда JWT токен в формате: Bearer {token}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	database.Connect()

	if err := services.RecalculateStatistics(); err != nil {
		log.Printf("Ошибка пересчета статистики: %v", err)
	}

	go services.AutoDeleteOldAds()

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Total-Count, Range, Content-Range, Accept")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "43200")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "http://localhost:3001", "http://localhost:3002", "http://127.0.0.1:5173", "http://127.0.0.1:3000", "http://127.0.0.1:3001", "http://127.0.0.1:3002"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Total-Count", "Range", "Content-Range", "Accept"},
		ExposeHeaders:    []string{"X-Total-Count", "Content-Range"},
		AllowCredentials: true,
		AllowWildcard:    false,
		MaxAge:           12 * 3600,
	}

	router.Use(cors.New(config))
	router.Use(middleware.RequestTimeout(30 * time.Second))

	router.Static("/uploads", "./uploads")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/api/register", middleware.RateLimitRegister(), services.RegisterAccount)
	router.POST("/api/login", middleware.RateLimitLogin(), services.Login)
	router.POST("/api/refresh", services.RefreshAccessToken)
	router.POST("/api/logout", services.Logout)

	router.POST("/api/verify-email", middleware.RateLimitVerify(), services.VerifyEmail)
	router.POST("/api/resend-code", middleware.RateLimitVerify(), services.ResendVerificationCode)

	router.POST("/api/createnewads", handlers.CreateNewAds)
	router.GET("/api/ads", handlers.GetAdsByCategory)
	router.GET("/api/ads/random", handlers.GetRandomAds)
	router.GET("/api/listings/user/:nickname", handlers.GetAdsByNickname)
	router.GET("/api/getadcount", handlers.GetAdCount)
	router.POST("/api/ads/:id/view", handlers.IncrementAdViews)
	router.PUT("/api/ads/:id", middleware.AuthRequired(), handlers.UpdateAd)
	router.DELETE("/api/ads/:id", middleware.AuthRequired(), handlers.DeleteAd)
	router.POST("/api/reports", middleware.AuthRequired(), handlers.CreateReport)

	router.POST("/api/feedback", middleware.AuthRequired(), handlers.CreateFeedback)
	router.GET("/api/feedback/:nickname", handlers.GetFeedbacksByOwner)
	router.PUT("/api/feedback/:id/confirm", middleware.AuthRequired(), handlers.ConfirmFeedback)

	router.POST("/api/viewed-ads", middleware.AuthRequired(), handlers.AddViewedAd)
	router.GET("/api/viewed-ads", middleware.AuthRequired(), handlers.GetViewedAds)

	router.POST("/api/profile/update-background", middleware.AuthRequired(), handlers.UpdateProfileBackground)
	router.DELETE("/api/profile/delete-background", middleware.AuthRequired(), handlers.DeleteProfileBackground)

	router.POST("/api/profile/update-avatar", middleware.AuthRequired(), handlers.UpdateProfileAvatar)
	router.PUT("/api/profile/update-nickname", middleware.AuthRequired(), handlers.UpdateNickname)
	router.PUT("/api/profile/update-email", middleware.AuthRequired(), handlers.UpdateEmail)
	router.PUT("/api/profile/update-password", middleware.AuthRequired(), handlers.UpdatePassword)
	router.PUT("/api/profile/update-theme", middleware.AuthRequired(), handlers.UpdateTheme)
	router.PUT("/api/profile/update-description", middleware.AuthRequired(), handlers.UpdateDescription)
	router.PUT("/api/profile/update-telegram", middleware.AuthRequired(), handlers.UpdateTelegram)

	router.GET("/api/me", middleware.AuthRequired(), func(c *gin.Context) {
		nickname, exists := c.Get("nickname")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
			return
		}

		user, err := services.GetUserByNickname(nickname.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных пользователя"})
			return
		}

		var reviewsCount int64
		database.DB.Table("feedback_ads").
			Where("ad_owner_nickname = ? AND confirm_feedback = ?", nickname.(string), true).
			Count(&reviewsCount)

		services.UpdateLastSeen(nickname.(string))

		c.JSON(http.StatusOK, gin.H{
			"user_id":                   user.ID,
			"nickname":                  user.Nickname,
			"email":                     user.Email,
			"telegram":                  user.Telegram,
			"avatar":                    user.Avatar,
			"background_avatar_profile": user.BackgroundAvatarProfile,
			"rating":                    user.Rating,
			"reviews_count":             reviewsCount,
			"user_role":                 user.UserRole,
			"user_description":          user.UserDescription,
			"theme":                     user.Theme,
			"last_seen_at":              user.LastSeenAt,
		})
	})

	port := ":8080"
	fmt.Printf("Сервер запущен на http://localhost%s\n", port)
	log.Fatal(router.Run(port))
}
