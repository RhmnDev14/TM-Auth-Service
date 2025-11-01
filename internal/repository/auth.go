package repository

import (
	"auth-service/internal/dto"
	"auth-service/internal/entity"
	"auth-service/internal/helper"
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type AuthRepo interface {
	Register(ctx context.Context, param dto.Register) error
	Login(ctx context.Context, email string) (*entity.User, error)
}

type authRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *authRepo {
	return &authRepo{db: db}
}

func (a *authRepo) Register(ctx context.Context, param dto.Register) error {
	user := entity.RegisterToEntity(param)
	err := a.db.WithContext(ctx).Create(&user).Error

	if err != nil {
		helper.ErrorHandle(err)
		if strings.Contains(err.Error(), "duplicate key value") {
			return fmt.Errorf(helper.UniqEmail)
		}

		if errors.Is(err, gorm.ErrInvalidData) {
			return fmt.Errorf(helper.InvalidData)
		}

		return fmt.Errorf(helper.FailedSave)
	}

	return nil
}

func (a *authRepo) Login(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User

	err := a.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		helper.ErrorHandle(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(helper.InvalidLogin)
		}
		return nil, fmt.Errorf(helper.FailedGet)
	}

	return &user, nil
}
