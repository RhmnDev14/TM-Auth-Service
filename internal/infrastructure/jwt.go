package infrastructure

import (
	"auth-service/internal/config"
	"auth-service/internal/helper"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userID uint) (string, error)
	ValidateToken(ctx context.Context, tokenString string) (context.Context, *jwt.Token, error)
}

type jwtService struct {
	cfg config.Config
}

func NewJWTService(cfg config.Config) *jwtService {
	return &jwtService{
		cfg: cfg,
	}
}

func (j *jwtService) GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"iss":     j.cfg.IssuerName,
		"exp":     time.Now().Add(j.cfg.JwtExpiresTime).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.cfg.JwtSignatureKy))
}

func (j *jwtService) ValidateToken(ctx context.Context, tokenString string) (context.Context, *jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			helper.ErrorHandle(errors.New(helper.InvalidMethodSigning))
			return nil, errors.New(helper.InternalServerError)
		}
		return []byte(j.cfg.JwtSignatureKy), nil
	})

	if err != nil {
		return ctx, nil, err
	}

	if !token.Valid {
		return ctx, nil, errors.New(helper.InvalidToken)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		helper.ErrorHandle(errors.New("cannot parse jwt claims"))
		return ctx, nil, errors.New(helper.InternalServerError)
	}

	userID, ok := claims["user_id"].(uint)
	if !ok {
		helper.ErrorHandle(errors.New("user_id not found in token claims"))
		return ctx, nil, errors.New(helper.InternalServerError)
	}
	keyCtx := "user_id"
	ctx = context.WithValue(ctx, keyCtx, userID)

	return ctx, token, nil
}
