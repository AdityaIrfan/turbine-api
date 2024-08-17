package services

import (
	"net/http"
	"net/mail"
	"strings"

	contract "pln/AdityaIrfan/turbine-api/contracts"
	"pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"

	"github.com/labstack/echo/v4"
)

type userService struct {
	userRepo     contract.IUserRepository
	divisionRepo contract.IDivisionRepository
	roleRepo     contract.IRoleRepository
}

func NewUserService(userRepo contract.IUserRepository, divisionRepo contract.IDivisionRepository, roleRepo contract.IRoleRepository) contract.IUserService {
	return &userService{
		userRepo:     userRepo,
		divisionRepo: divisionRepo,
		roleRepo:     roleRepo,
	}
}

func (u *userService) CreateUserByAdmin(c echo.Context, in *models.UserCreateByAdminRequest) error {
	// check existing username
	if exist, err := u.userRepo.IsUsernameExist(in.Username); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if exist {
		return helpers.Response(c, http.StatusBadRequest, "username sudah digunakan")
	}

	// check existing email
	if exist, err := u.userRepo.IsEmailExist(in.Email); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if exist {
		return helpers.Response(c, http.StatusBadRequest, "email sudah digunakan")
	}

	// check division by id
	division, err := u.divisionRepo.GetByIdWithSelectedFields(in.DivisionId, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if division.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "divisi tidak ditemukan")
	}

	user := in.ToModel()
	if err := u.userRepo.Create(user, "Division"); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil membuat user baru", user.ToResponse())
}

func (u *userService) UpdateUserByAdmin(c echo.Context, in *models.UserUpdateByAdminRequest) error {
	user, err := u.userRepo.GetById(in.Id, "Division")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "user tidak ditemukan")
	}

	if !user.IsGeneralUser() {
		return helpers.Response(c, http.StatusBadRequest, "kamu tidak memliki izin untuk mengakses fitur ini")
	}

	var anyUpdated bool

	if in.Role != nil && user.Role != *in.Role {
		if !models.IsUserRoleAvailable(*in.Role) {
			return helpers.Response(c, http.StatusBadRequest, "role tidak tersedia")
		}

		anyUpdated = true
		user.Role = *in.Role
	}

	if in.Status != nil && user.Status != *in.Status {
		if !models.IsUserStatusExist(*in.Status) {
			return helpers.Response(c, http.StatusBadRequest, "status user tidak tersedia")
		}

		switch *in.Status {
		case models.UserStatus_Active:
			user.ActivatedBy = in.UpdatedBy
		case models.UserStatus_BlockedByAdmin:
			user.BlockedBy = in.UpdatedBy
		}

		anyUpdated = true
		user.Status = *in.Status
	}

	if in.RadiusStatus != nil && user.RadiusStatus != *in.RadiusStatus {
		user.RadiusStatus = *in.RadiusStatus
		anyUpdated = true
	}

	if in.DivisionId != nil && user.DivisionId != *in.DivisionId {
		division, err := u.divisionRepo.GetByIdWithSelectedFields(*in.DivisionId, "*")
		if err != nil {
			return helpers.ResponseUnprocessableEntity(c)
		} else if division.IsEmpty() {
			return helpers.Response(c, http.StatusBadRequest, "divisi tidak ditemukan")
		}

		anyUpdated = true
		user.DivisionId = *in.DivisionId
		user.Division = division
	}

	if anyUpdated {
		if err := u.userRepo.Update(user, "Division"); err != nil {
			return helpers.ResponseUnprocessableEntity(c)
		}
	}

	return helpers.Response(c, http.StatusOK, "berhasil mengubah user", user.ToResponse())
}

func (u *userService) GetDetailUserByAdmin(c echo.Context, in *models.UserGetDetailRequest) error {
	user, err := u.userRepo.GetById(in.Id, "Division")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "user tidak ditemukan")
	}

	if !user.IsGeneralUser() {
		return helpers.Response(c, http.StatusBadRequest, "kamu tidak memliki izin untuk mengakses fitur ini")
	}

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan user", user.ToResponse())
}

func (u *userService) DeleteUserByAdmin(c echo.Context, in *models.UserDeleteByAdminRequest) error {
	user, err := u.userRepo.GetByIdWithSelectedFields(in.Id, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		return helpers.ResponseForbiddenAccess(c)
	}

	if !user.IsGeneralUser() {
		return helpers.Response(c, http.StatusBadRequest, "kamu tidak memliki izin untuk mengakses fitur ini")
	}

	if err := u.userRepo.Delete(user); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil menghapus user")
}

func (u *userService) GetListUserWithPaginateByAdmin(c echo.Context, cursor *helpers.Cursor) error {
	users, pagination, err := u.userRepo.GetAllWithPaginate(cursor, models.UserRole_User)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	var userRes = []*models.UserListResponse{}
	for _, user := range users {
		userRes = append(userRes, user.ToResponseList())
	}

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan semua user", userRes, pagination)
}

func (u *userService) GenerateUserPasswordByAdmin(c echo.Context, in *models.GeneratePasswordByAdmin) error {
	user, err := u.userRepo.GetById(in.Id)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		return helpers.ResponseForbiddenAccess(c)
	}

	if !user.IsGeneralUser() {
		return helpers.Response(c, http.StatusBadRequest, "kamu tidak memliki izin untuk mengakses fitur ini")
	}

	password := helpers.GenerateRandomString(10)
	salt, hash, err := helpers.GenerateHashAndSalt(password)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	user.PasswordSalt = salt
	user.PasswordHash = hash

	if err := u.userRepo.Update(user); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil membuat password baru", map[string]interface{}{
		"Password": password,
	})
}

func (u *userService) UpdateMyProfile(c echo.Context, in *models.UserUpdateRequest) error {
	user, err := u.userRepo.GetById(in.Id, "Division")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "user tidak ditemukan")
	}

	var anyUpdated bool

	if in.Username != nil && *in.Username != "" && user.Username != *in.Username {
		userByUsername, err := u.userRepo.GetByUsernameWithSelectedFields(*in.Username, "id")
		if err != nil {
			return helpers.ResponseUnprocessableEntity(c)
		} else if !userByUsername.IsEmpty() && userByUsername.Id != user.Id {
			return helpers.Response(c, http.StatusBadRequest, "username sudah digunakan")
		}

		anyUpdated = true
		user.Username = *in.Username
	}

	if in.Email != nil && *in.Email != "" && user.Email != *in.Email {
		if _, err := mail.ParseAddress(*in.Email); err != nil {
			return helpers.Response(c, http.StatusBadRequest, "Email is not valid")
		}

		userByEmail, err := u.userRepo.GetByEmailWithSelectedFields(*in.Email, "id")
		if err != nil {
			return helpers.ResponseUnprocessableEntity(c)
		} else if !userByEmail.IsEmpty() && userByEmail.Id != user.Id {
			return helpers.Response(c, http.StatusBadRequest, "email sudah digunakan")
		}

		anyUpdated = true
		user.Email = *in.Email
	}

	if in.Name != nil && *in.Name != "" && user.Name != *in.Name {
		anyUpdated = true
		user.Name = *in.Name
	}

	if in.Phone != nil && *in.Phone != "" && user.Phone != *in.Phone {
		phone := *in.Phone
		if strings.Index(phone, "0") == 0 {
			phone = "62" + phone[1:]
		}
		in.Phone = &phone

		// check phone
		userByPhone, err := u.userRepo.GetByPhoneWithSelectedFields(*in.Phone, "id")
		if err != nil {
			return helpers.ResponseUnprocessableEntity(c)
		} else if !userByPhone.IsEmpty() && userByPhone.Id != user.Id {
			return helpers.Response(c, http.StatusBadRequest, "nomor telefon sudah digunakan")
		}
		user.Phone = *in.Phone
		anyUpdated = true
	}

	if !anyUpdated {
		if err := u.userRepo.Update(user, "Division"); err != nil {
			return helpers.ResponseUnprocessableEntity(c)
		}
	}

	return helpers.Response(c, http.StatusOK, "berhasil mengubah user", user.ToResponse())
}

func (u *userService) GetMyProfile(c echo.Context, id string) error {
	user, err := u.userRepo.GetById(id, "Division")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "user tidak ditemukan")
	}

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan user", user.ToResponse())
}

func (u *userService) ChangeMyPassword(c echo.Context, in *models.UserChangePasswordRequest) error {
	user, err := u.userRepo.GetById(in.Id)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		return helpers.ResponseForbiddenAccess(c)
	}

	salt, hash, err := helpers.GenerateHashAndSalt(in.Password)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	user.PasswordSalt = salt
	user.PasswordHash = hash

	return helpers.Response(c, http.StatusOK, "berhasil mengubah password")
}
