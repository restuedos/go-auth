package database

import (
	"log"

	"go-auth/internal/domain/entity"
	"go-auth/internal/domain/repository"

	"gorm.io/gorm"
)

func SeedAdminUser(db *gorm.DB, userRepo repository.UserRepository) {
	admins, err := userRepo.FindByRole(entity.RoleAdmin)
	if err == nil && len(admins) > 0 {
		log.Println("Admin already exists, skipping creation.")
		return
	}

	admin := entity.User{
		Email:    "admin@example.com",
		Password: "admin123", // Change this after first login
		Name:     "Default Admin",
		Role:     entity.RoleAdmin,
	}

	if err := admin.HashPassword(); err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}

	if err := userRepo.Create(&admin); err != nil {
		log.Fatalf("Failed to create default admin: %v", err)
	}

	log.Println("Default admin created successfully.")
}
