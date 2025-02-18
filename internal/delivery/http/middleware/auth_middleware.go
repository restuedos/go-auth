package middleware

import (
	"go-auth/internal/config"
	"go-auth/internal/delivery/http/handler/responses"
	"go-auth/internal/usecase"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var cfg *config.Config
var authUsecase usecase.AuthUsecase

func Init(c *config.Config, uc usecase.AuthUsecase) {
	cfg = c
	authUsecase = uc
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, responses.ErrorResponse{
				Code:  http.StatusForbidden,
				Error: "authorization header is required",
			})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		blacklisted, _ := authUsecase.IsTokenBlacklisted(tokenString)
		if blacklisted {
			c.JSON(http.StatusUnauthorized, responses.ErrorResponse{
				Code:  http.StatusForbidden,
				Error: "token has been blacklisted",
			})
			c.Abort()
			return
		}

		_, claims, err := ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, responses.ErrorResponse{
				Code:  http.StatusForbidden,
				Error: "invalid token",
			})
			c.Abort()
			return
		}

		c.Set("token", tokenString)
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Set("role", claims["role"].(string))
		c.Next()
	}
}

func GenerateToken(userID uint, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(cfg.JWTSecret))
}

func ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, err
	}

	return token, claims, nil
}
