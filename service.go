package main

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"

	database "github.com/ScMofeoluwa/minalytics/database/sqlc"
)

var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrUserNotCreated = errors.New("failed to create user")
)

type AnalyticsService struct {
	queries *database.Queries
}

func NewAnalyticsService(queries *database.Queries) *AnalyticsService {
	return &AnalyticsService{queries: queries}
}

func (s *AnalyticsService) SignIn(ctx context.Context, email string) (string, error) {
	userId, err := s.queries.FindOrCreateUser(ctx, email)
	if err != nil {
		return "", ErrUserNotCreated
	}
	return CreateJWT(userId.String())
}

func CreateJWT(userId string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(viper.GetString("TokenSecret"))
	return token.SignedString(secret)
}

func VerifyJWT(token string) (jwt.MapClaims, error) {
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(viper.GetString("TokenSecret")), nil
	}

	jwtToken, err := jwt.Parse(token, keyfunc)
	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
