package database

import (
	addcells "arizonagamesstore/backend/migrations/add_cells"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		errEnv = godotenv.Load("../.env")
		if errEnv != nil {
			log.Println("Warning: .env file not found, using environment variables")
		}
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("error connecting the database. Error: %s", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get database instance: %s", err)
	}

	sqlDB.SetMaxIdleConns(10)

	sqlDB.SetMaxOpenConns(100)

	sqlDB.SetConnMaxLifetime(time.Hour)

	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	log.Println("🚀 Database success the connected :3")
	log.Println("✅ Connection pool configured: MaxIdle=10, MaxOpen=100")

	RunMigrations()
}

func RunMigrations() {
	migrateURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	migrationsPath := "file://migrations"
	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		migrationsPath = "file://../migrations"
	}

	m, err := migrate.New(migrationsPath, migrateURL)
	if err != nil {
		log.Fatalf("❌ Failed to create migrate instance: %s", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("✅ No new migrations to apply")
		} else {
			log.Fatalf("❌ Migration failed: %s", err)
		}
	} else {
		log.Println("✅ Migrations applied successfully")
	}

	err = addcells.SeedStatistics(DB)
	if err != nil {
		log.Fatalf("❌ Failed to seed statistics data: %s", err)
	}
	log.Println("✅ Statistics data seeded successfully")

	CreateViewedAdsTable()

	CreateFeedbackAdsTable()
}

func CreateViewedAdsTable() {
	var exists bool
	err := DB.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'viewed_ads')").Scan(&exists).Error
	if err != nil {
		log.Printf("⚠️ Error checking viewed_ads table: %s", err)
		return
	}

	if exists {
		log.Println("✅ viewed_ads table already exists")
		return
	}

	sqlScript := `
		CREATE TABLE viewed_ads (
			id SERIAL PRIMARY KEY,
			user_nickname VARCHAR(255) NOT NULL,
			ad_id INTEGER NOT NULL,
			viewed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_user FOREIGN KEY (user_nickname) REFERENCES accounts(nickname) ON DELETE CASCADE,
			CONSTRAINT fk_viewed_ad FOREIGN KEY (ad_id) REFERENCES ads(id) ON DELETE CASCADE,
			CONSTRAINT unique_user_ad UNIQUE(user_nickname, ad_id)
		);

		CREATE INDEX idx_viewed_ads_user ON viewed_ads(user_nickname);
		CREATE INDEX idx_viewed_ads_ad ON viewed_ads(ad_id);
		CREATE INDEX idx_viewed_ads_time ON viewed_ads(viewed_at DESC);
	`

	if err := DB.Exec(sqlScript).Error; err != nil {
		log.Printf("❌ Failed to create viewed_ads table: %s", err)
	} else {
		log.Println("✅ viewed_ads table created successfully")
	}
}

func CreateFeedbackAdsTable() {
	var exists bool
	err := DB.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'feedback_ads')").Scan(&exists).Error
	if err != nil {
		log.Printf("⚠️ Error checking feedback_ads table: %s", err)
		return
	}

	if exists {
		log.Println("✅ feedback_ads table already exists")
		return
	}

	sqlScript := `
		CREATE TABLE feedback_ads (
			id SERIAL PRIMARY KEY,
			ad_id INTEGER NOT NULL,
			reviewer_nickname VARCHAR(255) NOT NULL,
			ad_owner_nickname VARCHAR(255) NOT NULL,
			rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
			review_text TEXT NOT NULL,
			proof_image VARCHAR(500) NOT NULL,
			confirm_feedback BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_ad FOREIGN KEY (ad_id) REFERENCES ads(id) ON DELETE CASCADE,
			CONSTRAINT fk_reviewer FOREIGN KEY (reviewer_nickname) REFERENCES accounts(nickname) ON DELETE CASCADE,
			CONSTRAINT fk_ad_owner FOREIGN KEY (ad_owner_nickname) REFERENCES accounts(nickname) ON DELETE CASCADE
		);

		CREATE INDEX idx_feedback_ads_ad_owner ON feedback_ads(ad_owner_nickname);
		CREATE INDEX idx_feedback_ads_reviewer ON feedback_ads(reviewer_nickname);
		CREATE INDEX idx_feedback_ads_confirm ON feedback_ads(confirm_feedback);
	`

	if err := DB.Exec(sqlScript).Error; err != nil {
		log.Printf("❌ Failed to create feedback_ads table: %s", err)
	} else {
		log.Println("✅ feedback_ads table created successfully")
	}
}
