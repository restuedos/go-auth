package requests

type RegisterRequest struct {
	Email    string `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password string `json:"password" gorm:"not null" validate:"required,min=8"`
	Name     string `json:"name"`
}
