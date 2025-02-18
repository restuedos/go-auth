package usecase

import (
	"fmt"
	"go-auth/internal/delivery/http/handler/requests"
	"go-auth/internal/domain/entity"
	"go-auth/internal/domain/repository"

	"gorm.io/gorm"
)

type UserUsecase interface {
	GetAllUsers() ([]entity.User, error)
	GetUser(userID uint) (*entity.User, error)
	UpdateUser(userID uint, req requests.UpdateUserRequest) error
	DeleteUser(userID uint) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) GetAllUsers() ([]entity.User, error) {
	return u.userRepo.GetAllUsers()
}

func (u *userUsecase) GetUser(userID uint) (*entity.User, error) {
	return u.userRepo.FindByID(userID)
}

func (u *userUsecase) UpdateUser(userID uint, req requests.UpdateUserRequest) error {
	updates := make(map[string]interface{})

	if req.Email != "" {
		existingUser, err := u.userRepo.FindByEmail(req.Email)
		if err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to check existing email: %w", err)
		}
		if existingUser != nil && existingUser.ID != userID {
			return fmt.Errorf("email already in use")
		}
		updates["email"] = req.Email
	}

	if req.Name != "" {
		updates["name"] = req.Name
	}

	if req.Password != "" {
		user := entity.User{Password: req.Password}
		if err := user.HashPassword(); err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		updates["password"] = user.Password
	}

	if len(updates) == 0 {
		return fmt.Errorf("no valid fields to update")
	}

	return u.userRepo.Update(userID, updates)
}

func (u *userUsecase) DeleteUser(userID uint) error {
	return u.userRepo.Delete(userID)
}
