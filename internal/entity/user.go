package entity

import (
	"auth-service/internal/dto"
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Name      string    `json:"name" gorm:"not null"`
	Password  string    `json:"password" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func RegisterToEntity(param dto.Register) User {
	return User{
		Name:      param.Name,
		Password:  param.Password,
		Email:     param.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
