package requests

type CreateUserRequest struct {
	Email    string `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password string `json:"password" gorm:"not null" validate:"required,min=8"`
	Name     string `json:"name" gorm:"not null" validate:"required"`
}

type UpdateUserRequest struct {
	Email    string `json:"email" gorm:"unique;not null" validate:"optional,email"`
	Password string `json:"password" gorm:"not null" validate:"optional,min=8"`
	Name     string `json:"name" gorm:"not null" validate:"optional"`
}
