package models

import (
	"time"
	"turbine-api/helpers"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type UserStatus uint8
type UserRole uint8

const (
	UserStatus_InActive       UserStatus = 0
	UserStatus_Active         UserStatus = 1
	UserStatus_BlockedByAdmin UserStatus = 2

	UserRole_SuperAdmin UserRole = 1
	UserRole_Admin      UserRole = 2
	UserRole_User       UserRole = 3
)

type User struct {
	Id         string `gorm:"column:id"`
	Name       string `gorm:"column:name"`
	Username   string `gorm:"column:username"`
	Email      string `gorm:"column:email"`
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

func (u *User) IsEmpty() bool {
	return u == nil
}

func (u *User) IsActive() bool {
	return u.Status == UserStatus_Active
}

func (u *User) IsInActive() bool {
	return u.Status == UserStatus_InActive
}

func (u *User) IsBlockedByAdmin() bool {
	return u.Status == UserStatus_BlockedByAdmin
}

func (u *User) IsAdmin() bool {
	return u.Role == UserRole_Admin
}

func (u *User) IsSuperAdmin() bool {
	return u.Role == UserRole_SuperAdmin
}

func (u *User) ToResponse() *UserResponse {
	res := &UserResponse{
		Name:      u.Name,
		Username:  u.Username,
		Email:     u.Email,
		Division:  string(u.Division.Type),
		Role:      u.GetUserRoleInString(),
		Status:    u.GetUserStatusInString(),
		CreatedAt: u.CreatedAt.Format(helpers.DefaultTimeFormat),
	}

	if u.UpdatedAt != nil {
		res.UpdatedAt = u.UpdatedAt.Format(helpers.DefaultTimeFormat)
	}

	return res
}

type UserResponse struct {
	Name      string `json:"Name"`
	Username  string `json:"Username"`
	Email     string `json:"Email"`
	Division  string `json:"Division"`
	Role      string `json:"Role"`
	Status    string `json:"Status"`
	CreatedAt string `json:"CreatedAt"`
	UpdatedAt string `json:"UpdatedAt"`
}

func (u *User) ToResponseList() *UserListResponse {
	return &UserListResponse{
		Name:     u.Name,
		Division: string(u.Division.Type),
		Status:   u.GetUserStatusInString(),
	}
}

type UserListResponse struct {
	Name     string `json:"Name"`
	Division string `json:"Division"`
	Status   string `json:"Status"`
}

type UserAdminCreateByAdminRequest struct {
	Name       string `json:"Name"`
	Username   string `json:"useranme"`
	Email      string `json:"Email"`
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
	Role       *UserRole   `json:"Role"`
	DivisionId *string     `json:"DivisionId"`
	Status     *UserStatus `json:"UserStatus"`
}

type UserUpdateRequest struct {
	Id       string
	Name     *string `json:"Name"`
	Username *string `json:"Username"`
	Email    *string `json:"Email"`
}

type UserChangePasswordRequest struct {
	Id                   string
	Password             string `json:"Password" validate:"required"`
	PasswordConfirmation string `json:"PasswordConfirmation" validate:"required,eqfield:Password"`
}

type GeneratePasswordByAdmin struct {
	Id      string
	AdminId string
}

type UserDeleteByAdminRequest struct {
	Id      string
	AdminId string
}

func IsUserRoleAvailable(userRole UserRole) bool {
	switch userRole {
	case UserRole_Admin:
		return true
	case UserRole_User:
		return true
	default:
		return false
	}
}

func IsUserStatusExist(userStatus UserStatus) bool {
	switch userStatus {
	case UserStatus_Active:
		return true
	case UserStatus_InActive:
		return true
	default:
		return false
	}
}

func (u *User) GetUserRoleInString() string {
	switch u.Role {
	case UserRole_SuperAdmin:
		return "super admin"
	case UserRole_Admin:
		return "admin"
	case UserRole_User:
		return "user"
	default:
		return ""
	}
}

func (u *User) GetUserStatusInString() string {
	switch u.Status {
	case UserStatus_Active:
		return "active"
	case UserStatus_InActive:
		return "inactive"
	case UserStatus_BlockedByAdmin:
		return "blocked by admin"
	default:
		return ""
	}
}

type UserGetDetailRequest struct {
	Id string
}
