package services

import (
	"net/http"

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

func (u *userService) CreateUserAdminByAdmin(c echo.Context, in *models.UserAdminCreateByAdminRequest) error {
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
	if err := u.userRepo.Create(user); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil membuat user baru", user.ToResponse())
}

func (u *userService) UpdateByAdmin(c echo.Context, in *models.UserUpdateByAdminRequest) error {
	user, err := u.userRepo.GetById(in.Id, "Division")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "user tidak ditemukan")
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

		anyUpdated = true
		user.Status = *in.Status
	}

	if in.DivisionId != nil && user.DivisionId != *in.DivisionId {
		division, err := u.divisionRepo.GetByIdWithSelectedFields(*in.DivisionId, "id")
		if err != nil {
			return helpers.ResponseUnprocessableEntity(c)
		} else if division.IsEmpty() {
			return helpers.Response(c, http.StatusBadRequest, "divisi tidak ditemukan")
		}

		anyUpdated = true
		user.DivisionId = *in.DivisionId
	}

	if anyUpdated {
		if err := u.userRepo.Update(user, "Division"); err != nil {
			return helpers.ResponseUnprocessableEntity(c)
		}
	}

	return helpers.Response(c, http.StatusOK, "berhasil mengubah user", user.ToResponse())
}

func (u *userService) Update(c echo.Context, in *models.UserUpdateRequest) error {
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

	if !anyUpdated {
		if err := u.userRepo.Update(user, "Division"); err != nil {
			return helpers.ResponseUnprocessableEntity(c)
		}
	}

	return helpers.Response(c, http.StatusOK, "berhasil mengubah user", user.ToResponse())
}

func (u *userService) GetDetailByAdmin(c echo.Context, in *models.UserGetDetailRequest) error {
	user, err := u.userRepo.GetById(in.Id, "Division")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "user tidak ditemukan")
	}

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan user", user.ToResponse())
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

func (u *userService) DeleteByAdmin(c echo.Context, in *models.UserDeleteByAdminRequest) error {
	user, err := u.userRepo.GetByIdWithSelectedFields(in.Id, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		return helpers.ResponseForbiddenAccess(c)
	}

	if err := u.userRepo.Delete(user); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil menghaous user")
}

func (u *userService) GetListWithPaginateByAdmin(c echo.Context, cursor *helpers.Cursor) error {
	users, pagination, err := u.userRepo.GetAllWithPaginate(cursor)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	var userRes []*models.UserListResponse
	for _, user := range users {
		userRes = append(userRes, user.ToResponseList())
	}

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan semua user", userRes, pagination)
}

func (u *userService) ChangePassword(c echo.Context, in *models.UserChangePasswordRequest) error {
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

	if err := u.userRepo.Update(user); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil mengubah password")
}

func (u *userService) GeneratePasswordByAdmin(c echo.Context, in *models.GeneratePasswordByAdmin) error {
	user, err := u.userRepo.GetById(in.Id)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		return helpers.ResponseForbiddenAccess(c)
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
