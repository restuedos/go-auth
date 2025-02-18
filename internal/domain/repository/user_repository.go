package repository

import (
	"go-auth/internal/domain/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetAllUsers() ([]entity.User, error)
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByRole(role entity.Role) ([]entity.User, error)
	FindByID(id uint) (*entity.User, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAllUsers() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByRole(role entity.Role) ([]entity.User, error) {
	var users []entity.User
	err := r.db.Where("role = ?", role).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&entity.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}
