package responses

import "go-auth/internal/domain/entity"

type UserResponse struct {
	Users []entity.User `json:"users"`
}
