package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type UserStatus uint8
type UserRole uint8

const (
	UserStatus_InActive UserStatus = 0
	UserStatus_Active   UserStatus = 1
	UserStatus_Block    UserStatus = 2

	UserRole_Admin UserRole = 1
	UserRole_User  UserRole = 2
)

type User struct {
	Id         string `gorm:"column:id"`
	Name       string `gorm:"column:name"`
	Username   string `gorm:"column:username"`
	DivisionId string `gorm:"column:division_id"`
	// RoleId       string          `gorm:"column:role_id"`
	Role         UserRole   `gorm:"column:role"`
	Status       UserStatus `gorm:"status"`
	PasswordHash string     `gorm:"column:password_hash"`
	PasswordSalt string     `gorm:"column:password_salt"`
	CreatedAt    *time.Time `gorm:"column:created_at"`
	// CreatedBy    string          `gorm:"column:created_by"`
	UpdatedAt *time.Time `gorm:"column:updated_at;<-:update"`
	// UpdatedBy    *string         `gorm:"column:updated_by"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at"`
	// DeletedBy    *string         `gorm:"column:deleted_by"`

	// Role     *Role     `gorm:"foreignKey:RoleId;references:Id"`
	Division *Division `gorm:"foreignKey:DivisionId;references:Id"`
}

type UserAdminCreateByAdminRequest struct {
	Name       string `json:"Name"`
	Username   string `json:"useranme"`
	DivisionId string `json:"DivisionId"`
}

func (u *UserAdminCreateByAdminRequest) ToModel() *User {
	id := ulid.Make().String()

	return &User{
		Id:           id,
		Name:         u.Name,
		Username:     u.Username,
		DivisionId:   u.DivisionId,
		Role:         UserRole_Admin,
		Status:       UserStatus_Active,
		PasswordHash: "",
		PasswordSalt: "",
	}
}

type UserUpdateByAdminRequest struct {
	Id         string
	Role       *UserRole `json:"Role"`
	DivisionId *string   `json:"DivisionId"`
	AdminId    string
}

type UserUpdateRequest struct {
	Id       string
	Name     string `json:"Name"`
	Username string `json:"Username"`
}

type UserChangePassword struct {
	Password             string `json:"Password" validate:"required"`
	PasswordConfirmation string `json:"PasswordConfirmation" validate:"required,eqfield:Password"`
}

type UserDeleteByAdminRequest struct {
	Id      string
	AdminId string
}
