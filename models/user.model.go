package models

import (
	"strings"
	"time"

	"pln/AdityaIrfan/turbine-api/helpers"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

var UserDefaultSort = map[string]string{
	"Name":      "name",
	"Username":  "Username",
	"CreatedAt": "createdat",
}

var UserDefaultFilter = map[string]string{
	"Status": "status",
}

var UserDefaultFilterBySuperAdmin = map[string]string{
	"Status": "status",
	"Role":   "role",
}

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
	Id           string          `gorm:"column:id"`
	Name         string          `gorm:"column:name"`
	Username     string          `gorm:"column:username"`
	Email        string          `gorm:"column:email"`
	Phone        string          `gorm:"column:phone"`
	DivisionId   string          `gorm:"column:division_id"`
	Role         UserRole        `gorm:"column:role"`
	Status       UserStatus      `gorm:"column:status"`
	RadiusStatus bool            `gorm:"column:radius_status"`
	ActivatedBy  string          `gorm:"column:activated_by"`
	BlockedBy    string          `gorm:"column:blocked_by"`
	PasswordHash string          `gorm:"column:password_hash"`
	PasswordSalt string          `gorm:"column:password_salt"`
	CreatedAt    *time.Time      `gorm:"column:created_at"`
	CreatedBy    string          `gorm:"column:created_by"`
	UpdatedAt    *time.Time      `gorm:"column:updated_at;<-:update"`
	DeletedAt    *gorm.DeletedAt `gorm:"column:deleted_at"`
	DeletedBy    string          `gorm:"column:deleted_by"`

	CreatedByUser *User     `gorm:"foreignKey:CreatedBy;references:Id"`
	Division      *Division `gorm:"foreignKey:DivisionId;references:Id"`
}

func (u *User) TableName() string {
	return "users"
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

func (u *User) IsGeneralUser() bool {
	return u.Role == UserRole_User
}

func (u *User) ToResponse() *UserResponse {
	res := &UserResponse{
		Id:           u.Id,
		Name:         u.Name,
		Username:     u.Username,
		Email:        u.Email,
		Division:     string(u.Division.Name),
		Role:         u.GetUserRoleInString(),
		Status:       u.GetUserStatusInString(),
		RadiusStatus: u.RadiusStatus,
		CreatedAt:    u.CreatedAt.Format(helpers.DefaultTimeFormat),
	}

	if u.UpdatedAt != nil {
		res.UpdatedAt = u.UpdatedAt.Format(helpers.DefaultTimeFormat)
	}

	if strings.Index(u.Phone, "62") == 0 {
		res.Phone = "+" + u.Phone
	}

	return res
}

type UserResponse struct {
	Id           string `json:"Id"`
	Name         string `json:"Name"`
	Username     string `json:"Username"`
	Email        string `json:"Email"`
	Phone        string `json:"Phone"`
	Division     string `json:"Division"`
	Role         string `json:"Role"`
	Status       string `json:"Status"`
	RadiusStatus bool   `json:"RadiusStatus"`
	CreatedAt    string `json:"CreatedAt"`
	UpdatedAt    string `json:"UpdatedAt"`
}

func (u *User) ToResponseList() *UserListResponse {
	return &UserListResponse{
		Id:       u.Id,
		Name:     u.Name,
		Division: string(u.Division.Name),
		Status:   u.GetUserStatusInString(),
	}
}

type UserListResponse struct {
	Id       string `json:"Id"`
	Name     string `json:"Name"`
	Division string `json:"Division"`
	Status   string `json:"Status"`
}

type UserCreateByAdminRequest struct {
	Name       string `json:"Name" form:"Name" validate:"required"`
	Username   string `json:"useranme" form:"Username" validate:"required"`
	Email      string `json:"Email" form:"Email" validate:"required,email"`
	DivisionId string `json:"DivisionId" form:"DivisionId" validate:"required"`
	Phone      string `json:"Phone" form:"Phone" validate:"required"`
	Password   string `json:"Password" form:"Password" validate:"required"`
	CreatedBy  string
}

func (u *UserCreateByAdminRequest) ToModel() *User {
	id := ulid.Make().String()

	return &User{
		Id:           id,
		Name:         u.Name,
		Username:     u.Username,
		Phone:        u.Phone,
		Email:        u.Email,
		DivisionId:   u.DivisionId,
		Role:         UserRole_User,
		Status:       UserStatus_Active,
		RadiusStatus: true,
		PasswordHash: "",
		PasswordSalt: "",
		CreatedBy:    u.CreatedBy,
	}
}

type UserUpdateByAdminRequest struct {
	Id           string
	Role         *UserRole   `json:"Role" form:"Role"`
	DivisionId   *string     `json:"DivisionId" form:"DivisionId"`
	Status       *UserStatus `json:"Status" form:"Status"`
	RadiusStatus *bool       `json:"RadiusStatus" form:"RadiusStatus"`
	UpdatedBy    string
}

type UserUpdateRequest struct {
	Id       string
	Name     *string `json:"Name" form:"Name"`
	Username *string `json:"Username" form:"Username"`
	Email    *string `json:"Email" form:"Email"`
	Phone    *string `json:"Phone" form:"Phone"`
}

type UserChangePasswordRequest struct {
	Id                   string
	Password             string `json:"Password" form:"Password" validate:"required"`
	PasswordConfirmation string `json:"PasswordConfirmation" form:"PasswordConfirmation" validate:"required"`
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
	case UserStatus_BlockedByAdmin:
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

func (u *User) IsRadiusStatusActive() bool {
	return u.RadiusStatus
}

func (u *User) GetSource() string {
	switch u.Role {
	case 1:
		return "main"
	case 2:
		return "admin"
	default:
		return "user"
	}
}

// func (u *User) GeneratePassword(password string) error {
// 	// Generate a 32-byte key for AES-256
// 	key, err := helpers.GenerateKey(32)
// 	if err != nil {
// 		log.Error().Err(errors.New("ERROR GENERATING KEY : " + err.Error())).Msg("")
// 		return err
// 	}

// 	// Encrypt the plaintext
// 	cipherText, iv, err := helpers.Encrypt([]byte(password), key)
// 	if err != nil {
// 		log.Error().Err(errors.New("ERROR PASSWORD ENCRYPTION : " + err.Error())).Msg("")
// 		return err
// 	}

// 	u.Password = string(cipherText)
// 	u.Key = string(key)
// 	u.Iv = string(iv)
// 	return nil
// }

// func (u *User) GetPassword() string {
// 	// Decrypt the ciphertext
// 	decryptedText, err := helpers.Decrypt([]byte(u.Password), []byte(u.Key), []byte(u.Iv))
// 	if err != nil {
// 		log.Error().Err(errors.New("ERROR DECRYPTING PASSWORD : " + err.Error()))
// 		return ""
// 	}

// 	return string(decryptedText)
// }

type UserCreateBySuperAdminRequest struct {
	Name       string   `json:"Name" form:"Name" validate:"required"`
	Username   string   `json:"useranme" form:"Username" validate:"required"`
	Email      string   `json:"Email" form:"Email" validate:"required,email"`
	DivisionId string   `json:"DivisionId" form:"DivisionId" validate:"required"`
	Role       UserRole `json:"Role" form:"Role" validate:"required"`
	Phone      string   `json:"Phone" form:"Phone" validate:"required"`
	Password   string   `json:"Password" form:"Password" validate:"required"`
	CreatedBy  string
}

func (u *UserCreateBySuperAdminRequest) ToModel() *User {
	id := ulid.Make().String()

	return &User{
		Id:           id,
		Name:         u.Name,
		Username:     u.Username,
		Phone:        u.Phone,
		Email:        u.Email,
		DivisionId:   u.DivisionId,
		Role:         u.Role,
		Status:       UserStatus_Active,
		RadiusStatus: true,
		PasswordHash: "",
		PasswordSalt: "",
		CreatedBy:    u.CreatedBy,
	}
}

type UserUpdateBySuperAdminRequest struct {
	Id           string
	Role         *UserRole   `json:"Role" form:"Role"`
	DivisionId   *string     `json:"DivisionId" form:"DivisionId"`
	Status       *UserStatus `json:"Status" form:"Status"`
	RadiusStatus *bool       `json:"RadiusStatus" form:"RadiusStatus"`
	UpdatedBy    string
}

type UserDeleteBySuperAdminRequest struct {
	Id      string
	AdminId string
}

type GeneratePasswordBySuperAdmin struct {
	Id      string
	AdminId string
}
