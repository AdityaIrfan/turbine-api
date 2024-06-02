package models

import (
	"time"
	"turbine-api/helpers"

	"github.com/oklog/ulid/v2"
)

type Register struct {
	Name                 string `json:"Name" validate:"required"`
	Username             string `json:"Username" validation:"required"`
	DivisionId           string `json:"DivisionId" validate:"required"`
	Password             string `json:"Password" validation:"required"`
	PasswordConfirmation string `json:"PasswordConfirmation" validate:"required,eqfield:Password"`
}

func (r *Register) ToModel() (*User, error) {
	id := ulid.Make().String()

	salt, hash, err := helpers.GenerateHashAndSalt(r.Password)
	if err != nil {
		return nil, err
	}

	return &User{
		Id:           id,
		Name:         r.Name,
		Username:     r.Username,
		DivisionId:   r.DivisionId,
		Role:         UserRole_User,
		Status:       UserStatus_InActive,
		PasswordHash: hash,
		PasswordSalt: salt,
	}, nil
}

type Login struct {
	Username string `json:"Username" validate:"required"`
	Password string `json:"Password" validate:"required"`
}

type AuthResponse struct {
	Name         string `json:"Name"`
	Division     string `json:"Division"`
	Token        string `json:"Token"`
	RefreshToken string `json:"RefreshToken"`
}

type RefreshToken struct {
	RefreshToken string `json:"RefreshToken" vaidate:"required"`
}

type RefreshTokenRedis struct {
	RefreshToken string `json:"refresh_token"`
	Exp          int64  `json:"exp"`
}

func (r *RefreshTokenRedis) IsActive() bool {
	return time.Now().After(time.Unix(r.Exp, 0))
}
