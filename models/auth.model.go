package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"pln/AdityaIrfan/turbine-api/helpers"
)

type Register struct {
	Name                 string `json:"Name" validate:"required"`
	Username             string `json:"Username" validation:"required"`
	DivisionId           string `json:"DivisionId" validate:"required"`
	Password             string `json:"Password" validation:"required"`
	PasswordConfirmation string `json:"PasswordConfirmation" validate:"required,eqfield=Password"`
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

type RefreshTokenRequest struct {
	RefreshToken string `json:"RefreshToken" validate:"required"`
}

type RefreshTokenRedis struct {
	RefreshToken string `json:"refresh_token"`
	Exp          int64  `json:"exp"`
	Active       int64  `json:"active"`
}

func (r *RefreshTokenRedis) IsActive() bool {
	return time.Now().After(time.Unix(r.Active, 0))
}

type Logout struct {
	Token string `json:"Token" validate:"required"`
}
