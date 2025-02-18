package handler

import (
	"go-auth/internal/delivery/http/handler/requests"
	"go-auth/internal/delivery/http/handler/responses"
	"go-auth/internal/delivery/http/middleware"
	"go-auth/internal/domain/entity"
	"go-auth/internal/usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body requests.RegisterRequest true "User data"
// @Success 201 {object} responses.TokenResponse "User registered successfully"
// @Failure 400 {object} responses.ErrorResponse "Invalid request data"
// @Failure 500 {object} responses.ErrorResponse "Failed to generate token"
// @Router /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var registerRequest requests.RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	user := entity.User{
		Email:    registerRequest.Email,
		Password: registerRequest.Password,
		Name:     registerRequest.Name,
	}

	if err := h.authUsecase.Register(&user); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	token, err := middleware.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "failed to generate token",
		})
		return
	}

	_ = h.authUsecase.CacheToken(user.ID, token)

	c.JSON(http.StatusCreated, responses.TokenResponse{
		Token: token,
	})
}

// Login godoc
// @Summary User login
// @Description Authenticates a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body requests.LoginRequest true "User login credentials"
// @Success 200 {object} responses.TokenResponse "Login successful"
// @Failure 400 {object} responses.ErrorResponse "Invalid request data"
// @Failure 401 {object} responses.ErrorResponse "Invalid credentials"
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var loginRequest requests.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	user, err := h.authUsecase.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse{
			Code:  http.StatusUnauthorized,
			Error: err.Error(),
		})
		return
	}

	cachedToken, err := h.authUsecase.GetCachedToken(user.ID)
	if err == nil && cachedToken != "" {
		c.JSON(http.StatusOK, responses.TokenResponse{
			Token: cachedToken,
		})
		return
	}

	token, err := middleware.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "failed to generate token",
		})
		return
	}

	_ = h.authUsecase.CacheToken(user.ID, token)

	c.JSON(http.StatusOK, responses.TokenResponse{
		Token: token,
	})
}

// Logout godoc
// @Summary Logout a user
// @Description Logs out the user by invalidating their token
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.SuccessResponse "Logout successful"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Failed to logout"
// @Router /logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	tokenString, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse{
			Code:  http.StatusUnauthorized,
			Error: "Unauthorized",
		})
		return
	}

	_, claims, err := middleware.ParseToken(tokenString.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse{
			Code:  http.StatusUnauthorized,
			Error: "Invalid token",
		})
		return
	}
	exp := time.Unix(int64(claims["exp"].(float64)), 0)

	err = h.authUsecase.BlacklistToken(tokenString.(string), exp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "Failed to blacklist token",
		})
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse{
		Message: "Successfully logged out",
	})
}
