package database

import (
	"fmt"
	"go-auth/internal/config"
	"go-auth/internal/domain/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
