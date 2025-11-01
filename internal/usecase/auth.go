package usecase

import (
	"auth-service/internal/dto"
	"auth-service/internal/helper"
	"auth-service/internal/infrastructure"
	"auth-service/internal/repository"
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthUc interface {
	Register(ctx context.Context, param dto.Register) error
	Login(ctx context.Context, param dto.Login) (dto.LoginResponse, error)
}

type authUc struct {
	repo repository.AuthRepo
	jwt  infrastructure.JWTService
}

func NewAuthUc(repo repository.AuthRepo, jwt infrastructure.JWTService) *authUc {
	return &authUc{
		repo: repo,
		jwt:  jwt,
	}
}

func (a *authUc) Register(ctx context.Context, param dto.Register) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)
	if err != nil {
		helper.ErrorHandle(fmt.Errorf(helper.GagalHas, err))
		return fmt.Errorf(helper.InternalServerError)
	}
	param.Password = string(bytes)
	return a.repo.Register(ctx, param)
}

func (a *authUc) Login(ctx context.Context, param dto.Login) (dto.LoginResponse, error) {
	user, err := a.repo.Login(ctx, param.Email)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(param.Password))
	if err != nil {
		return dto.LoginResponse{}, errors.New(helper.InvalidLogin)
	}

	token, err := a.jwt.GenerateToken(user.ID)
	if err != nil {
		helper.ErrorHandle(fmt.Errorf("failed to generate token: %v", err))
		return dto.LoginResponse{}, fmt.Errorf(helper.InternalServerError)
	}

	return dto.LoginResponse{
		Token: token,
		User: dto.UserData{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
