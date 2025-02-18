package requests

type LoginRequest struct {
	Email    string `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password string `json:"password" gorm:"not null" validate:"required"`
}
