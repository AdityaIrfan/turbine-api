package models

import (
	"time"

	"pln/AdityaIrfan/turbine-api/helpers"

	"github.com/oklog/ulid/v2"
)

type Register struct {
	Name                 string `json:"Name" form:"Name" validate:"required"`
	Username             string `json:"Username" form:"Username" validate:"required"`
	Email                string `json:"Email" form:"Email" validate:"required"`
	DivisionId           string `json:"DivisionId" form:"DivisionId" validate:"required"`
	Password             string `json:"Password" form:"Password" validation:"required"`
	PasswordConfirmation string `json:"PasswordConfirmation" form:"PasswordConfirmation" validate:"required,eqfield=Password"`
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
		Email:        r.Email,
		DivisionId:   r.DivisionId,
		Role:         UserRole_User,
		Status:       UserStatus_InActive,
		PasswordHash: hash,
		PasswordSalt: salt,
	}, nil
}

type Login struct {
	Username string `json:"Username" form:"Username" validate:"required"`
	Password string `json:"Password" form:"Password" validate:"required"`
}

type AuthResponse struct {
	Name         string `json:"Name"`
	Division     string `json:"Division"`
	IsAdmin      bool   `json:"IsAdmin"`
	Token        string `json:"Token"`
	RefreshToken string `json:"RefreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"RefreshToken" form:"RefreshToken" validate:"required"`
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
	Token string `json:"Token" form:"Token" validate:"required"`
}
