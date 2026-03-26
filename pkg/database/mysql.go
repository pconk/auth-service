package database

import (
	"auth-service/internal/config"
	"auth-service/internal/entity"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(cfg.DB_DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto Migrate Entity
	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}
