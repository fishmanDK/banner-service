package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/fishmanDK/avito_test_task/internal/storage"
	"github.com/fishmanDK/avito_test_task/models"
)

type AuthService struct {
	storage *storage.Storage
}

func NewAuthService(storage *storage.Storage) *AuthService {
	return &AuthService{
		storage: storage,
	}
}

func (a *AuthService) Authentication(user models.User) (models.Token, error) {
	const op = "service.Authentication"

	user.Password = HashPassword(user.Password)

	userRole, err := a.storage.DB.GetUserRole(user)
	if err != nil {
		return models.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	accessToken, err := CreateAccessToken(userRole)
	if err != nil {
		return models.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.Token{
		AccessToken: accessToken,
	}, nil
}

func (a *AuthService) CreateUser(newUser models.NewUser) error {
	const op = "service.CreateUser"

	newUser.Password = HashPassword(newUser.Password)
	err := a.storage.DB.CreateUser(newUser)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func HashPassword(password string) string {
	data := []byte(password + salt)
	hashData := sha256.Sum256(data)
	hashString := hex.EncodeToString(hashData[:])

	return hashString
}
