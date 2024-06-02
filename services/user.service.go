package services

import (
	"net/http"
	contract "turbine-api/contracts"
	"turbine-api/helpers"
	"turbine-api/models"

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
	userAdmin, err := u.userRepo.GetByIdWithSelectedFields(in.AdminId, "id")
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if userAdmin.IsEmpty() {
		return helpers.Response(c, http.StatusForbidden, "who are you? you do not have permission for this access")
	}

	// check existing username
	exist, err := u.userRepo.IsUsernameExist(in.Username)
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if exist {
		return helpers.Response(c, http.StatusBadRequest, "username already in use")
	}

	// check division by id
	division, err := u.divisionRepo.GetByIdWithSelectedFields(in.DivisionId, "id")
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if division.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "division not found")
	}

	user := in.ToModel()
	if err := u.userRepo.Create(user); err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	}

	return helpers.Response(c, http.StatusOK, "success create user")
}

func (u *userService) UpdateByAdmin(c echo.Context, in *models.UserUpdateByAdminRequest) error {
	userAdmin, err := u.userRepo.GetByIdWithSelectedFields(in.AdminId, "id")
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if userAdmin.IsEmpty() {
		return helpers.Response(c, http.StatusForbidden, "who are you? you do not have permission for this access")
	}

	user, err := u.userRepo.GetById(in.Id, "Division")
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if user.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "user not found")
	}

	var anyUpdated bool

	if in.Role != nil && user.Role != *in.Role {
		if !models.IsUserRoleAvailable(*in.Role) {
			return helpers.Response(c, http.StatusBadRequest, "role value is not available")
		}

		anyUpdated = true
		user.Role = *in.Role
	}

	if in.Status != nil && user.Status != *in.Status {
		if !models.IsUserStatusExist(*in.Status) {
			return helpers.Response(c, http.StatusBadRequest, "status value is not available")
		}

		anyUpdated = true
		user.Status = *in.Status
	}

	if in.DivisionId != nil && user.DivisionId != *in.DivisionId {
		division, err := u.divisionRepo.GetByIdWithSelectedFields(*in.DivisionId, "id")
		if err != nil {
			return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
		} else if division.IsEmpty() {
			return helpers.Response(c, http.StatusBadRequest, "division not found")
		}

		anyUpdated = true
		user.DivisionId = *in.DivisionId
	}

	if anyUpdated {
		if err := u.userRepo.Update(user, "Division"); err != nil {
			return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
		}
	}

	return helpers.Response(c, http.StatusOK, "success update user", user.ToResponse())
}

func (u *userService) Update(c echo.Context, in *models.UserUpdateRequest) error {
	user, err := u.userRepo.GetById(in.Id, "Division")
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if user.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "user not found")
	}

	var anyUpdated bool

	if in.Username != nil && user.Username != *in.Username {
		userByUsername, err := u.userRepo.GetByUsernameWithSelectedFields(*in.Username, "id")
		if err != nil {
			return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
		} else if !userByUsername.IsEmpty() && userByUsername.Id != user.Id {
			return helpers.Response(c, http.StatusBadRequest, "username already in use")
		}

		anyUpdated = true
		user.Username = *in.Username
	}

	if in.Name != nil && user.Name != *in.Name {
		anyUpdated = true
		user.Name = *in.Name
	}

	if !anyUpdated {
		if err := u.userRepo.Update(user, "Division"); err != nil {
			return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
		}
	}

	return helpers.Response(c, http.StatusOK, "success update user", user.ToResponse())
}

func (u *userService) GetDetailByAdmin(c echo.Context, in *models.UserGetDetailRequest) error {
	userAdmin, err := u.userRepo.GetByIdWithSelectedFields(in.AdminId, "id")
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if userAdmin.IsEmpty() {
		return helpers.Response(c, http.StatusForbidden, "who are you? you do not have permission for this access")
	}

	user, err := u.userRepo.GetById(in.Id, "Division")
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if user.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "user not found")
	}

	return helpers.Response(c, http.StatusOK, "success get user", user.ToResponse())
}

func (u *userService) GetMyProfile(c echo.Context, id string) error {
	user, err := u.userRepo.GetById(id, "Division")
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if user.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "user not found")
	}

	return helpers.Response(c, http.StatusOK, "success get user", user.ToResponse())
}

func (u *userService) DeleteByAdmin(c echo.Context, in *models.UserDeleteByAdminRequest) error {
	userAdmin, err := u.userRepo.GetByIdWithSelectedFields(in.AdminId, "id")
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if userAdmin.IsEmpty() {
		return helpers.Response(c, http.StatusForbidden, "who are you? you do not have permission for this access")
	}

	user, err := u.userRepo.GetByIdWithSelectedFields(in.Id, "id")
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if user.IsEmpty() {
		return helpers.Response(c, http.StatusForbidden, "who are you? you do not have permission for this access")
	}

	if err := u.userRepo.Delete(user); err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	}

	return helpers.Response(c, http.StatusOK, "success delete user")
}

func (u *userService) GetListWithPaginateByAdmin(c echo.Context, adminId string, cursor *helpers.Cursor) error {
	userAdmin, err := u.userRepo.GetByIdWithSelectedFields(adminId, "id")
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if userAdmin.IsEmpty() {
		return helpers.Response(c, http.StatusForbidden, "who are you? you do not have permission for this access")
	}

	users, pagination, err := u.userRepo.GetAllWithPaginate(cursor)
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	}

	var userRes []*models.UserListResponse
	for _, user := range users {
		userRes = append(userRes, user.ToResponseList())
	}

	return helpers.Response(c, http.StatusOK, "success get user list", userRes, pagination)
}

func (u *userService) ChangePassword(c echo.Context, in *models.UserChangePasswordRequest) error {
	user, err := u.userRepo.GetById(in.Id)
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if user.IsEmpty() {
		return helpers.Response(c, http.StatusForbidden, "who are you? you do not have permission for this access")
	}

	salt, hash, err := helpers.GenerateHashAndSalt(in.Password)
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	}

	user.PasswordSalt = salt
	user.PasswordHash = hash

	if err := u.userRepo.Update(user); err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	}

	return helpers.Response(c, http.StatusOK, "success change password")
}

func (u *userService) GeneratePasswordByAdmin(c echo.Context, in *models.GeneratePasswordByAdmin) error {
	userAdmin, err := u.userRepo.GetByIdWithSelectedFields(in.AdminId, "id")
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if userAdmin.IsEmpty() {
		return helpers.Response(c, http.StatusForbidden, "who are you? you do not have permission for this access")
	}

	user, err := u.userRepo.GetById(in.Id)
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if user.IsEmpty() {
		return helpers.Response(c, http.StatusForbidden, "who are you? you do not have permission for this access")
	}

	password := helpers.GenerateRandomString(10)
	salt, hash, err := helpers.GenerateHashAndSalt(password)
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	}

	user.PasswordSalt = salt
	user.PasswordHash = hash

	if err := u.userRepo.Update(user); err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	}

	return helpers.Response(c, http.StatusOK, "success generate password", map[string]interface{}{
		"Password": password,
	})
}
