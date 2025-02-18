package usecase

import (
	"context"
	"errors"
	"go-auth/internal/domain/entity"
	"go-auth/internal/domain/repository"
	"go-auth/internal/infrastructure/cache"
	"strconv"
	"time"
)

type AuthUsecase interface {
	Register(user *entity.User) error
	Login(email, password string) (*entity.User, error)
	CacheToken(userID uint, token string) error
	GetCachedToken(userID uint) (string, error)
	BlacklistToken(token string, expiresAt time.Time) error
	IsTokenBlacklisted(token string) (bool, error)
}

type authUsecase struct {
	userRepo repository.UserRepository
	cache    *cache.RedisCache
}

func NewAuthUsecase(userRepo repository.UserRepository, cache *cache.RedisCache) AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
		cache:    cache,
	}
}

func (a *authUsecase) Register(user *entity.User) error {
	existingUser, err := a.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("email already exists")
	}

	if err := user.HashPassword(); err != nil {
		return err
	}

	return a.userRepo.Create(user)
}

func (a *authUsecase) Login(email, password string) (*entity.User, error) {
	user, err := a.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := user.ComparePassword(password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (a *authUsecase) CacheToken(userID uint, token string) error {
	ctx := context.Background()
	expiration := 24 * time.Hour // Token expires in 24 hours
	cacheKey := "auth_token:" + strconv.FormatUint(uint64(userID), 10)
	return a.cache.Set(ctx, cacheKey, token, expiration)
}

func (a *authUsecase) GetCachedToken(userID uint) (string, error) {
	ctx := context.Background()
	cacheKey := "auth_token:" + strconv.FormatUint(uint64(userID), 10)
	token, err := a.cache.Get(ctx, cacheKey)
	if err != nil {
		return "", errors.New("token not found in cache")
	}
	return token, nil
}

func (a *authUsecase) BlacklistToken(token string, expiresAt time.Time) error {
	ctx := context.Background()
	cacheKey := "blacklist_token:" + token
	expiration := time.Until(expiresAt)

	if expiration <= 0 {
		expiration = time.Hour * 24 // Default Token expires in 24 hours if expiration is invalid
	}

	return a.cache.Set(ctx, cacheKey, "blacklisted", expiration)
}

func (a *authUsecase) IsTokenBlacklisted(token string) (bool, error) {
	ctx := context.Background()
	cacheKey := "blacklist_token:" + token
	_, err := a.cache.Get(ctx, cacheKey)
	if err != nil {
		return false, nil
	}
	return true, nil
}
