package handler

import (
	"go-auth/internal/delivery/http/handler/requests"
	"go-auth/internal/delivery/http/handler/responses"
	"go-auth/internal/domain/entity"
	"go-auth/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

// GetAllUsers godoc
// @Summary Get all users (Admin Only)
// @Description Retrieves a list of all registered users. Only accessible by admin users.
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {object} responses.UserResponse "List of users retrieved successfully"
// @Failure 403 {object} responses.ErrorResponse "Access denied"
// @Failure 500 {object} responses.ErrorResponse "Failed to fetch users"
// @Router /admin/users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	userRole, exists := c.Get("role")
	if !exists || userRole.(string) != string(entity.RoleAdmin) {
		c.JSON(http.StatusForbidden, responses.ErrorResponse{
			Code:  http.StatusForbidden,
			Error: "access denied",
		})
		return
	}

	users, err := h.userUsecase.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "failed to fetch users",
		})
		return
	}

	c.JSON(http.StatusOK, responses.UserResponse{Users: users})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Retrieves the authenticated user's profile
// @Tags user
// @Security BearerAuth
// @Produce json
// @Success 200 {object} entity.User "User profile retrieved successfully"
// @Failure 404 {object} responses.ErrorResponse "User not found"
// @Router /user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.userUsecase.GetUser(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Code:  http.StatusNotFound,
			Error: "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Updates the authenticated user's profile
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user body requests.UpdateUserRequest true "Updated user data"
// @Success 200 {object} responses.SuccessResponse "Profile updated successfully"
// @Failure 400 {object} responses.ErrorResponse "Invalid request data"
// @Failure 500 {object} responses.ErrorResponse "Failed to update profile"
// @Router /user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse{
			Code:  http.StatusUnauthorized,
			Error: "unauthorized",
		})
		return
	}

	var updates requests.UpdateUserRequest
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	if err := h.userUsecase.UpdateUser(userID.(uint), updates); err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse{Message: "profile updated successfully"})
}

// DeleteProfile godoc
// @Summary Delete user profile
// @Description Deletes the authenticated user's account
// @Tags user
// @Security BearerAuth
// @Success 200 {object} responses.SuccessResponse "Profile deleted successfully"
// @Failure 500 {object} responses.ErrorResponse "Failed to delete profile"
// @Router /user/profile [delete]
func (h *UserHandler) DeleteProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if err := h.userUsecase.DeleteUser(userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse{Message: "profile deleted successfully"})
}
