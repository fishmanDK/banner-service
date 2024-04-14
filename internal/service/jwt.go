package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

const (
	access_tokenTtl = time.Minute * 60
	signInKey       = "@(#tf53$*#$(RHfverib}#Rfrte)"
	salt            = "lsd2#tfv%2"
)

type JWT interface {
	CreateAccessToken(id int64) (string, error)
	CreateRefreshToken() (string, error)
	ParseToken(accessToken string) (*ParseDataUser, error)
}

type TokenClaims struct {
	jwt.StandardClaims
	Id    int64  `json:"id"`
	Role  string `json:"role"`
	Email string `json:"email"`
}

type ParseDataUser struct {
	Role string
}

func CreateAccessToken(role string) (string, error) {
	const op = "service.CreateAccessToken"

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(access_tokenTtl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Role: role,
	})

	signedAccessToken, err := accessToken.SignedString([]byte(signInKey))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return signedAccessToken, nil
}

func (s *AuthService) ParseToken(accessToken string) (string, error) {
	const op = "service.ParseToken"

	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: %w", op, errors.New("token verification error"))
		}
		return []byte(signInKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return "", errors.New("failed to parse token claims")
	}

	if !ok {
		return "", err
	}

	return claims.Role, nil
}
