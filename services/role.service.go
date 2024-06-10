package services

import (
	"net/http"
	contract "turbine-api/contracts"
	"turbine-api/helpers"
	"turbine-api/models"

	"github.com/labstack/echo/v4"
)

type roleService struct {
	roleRepo contract.IRoleRepository
}

func NewRoleService(roleRepo contract.IRoleRepository) contract.IRoleService {
	return &roleService{
		roleRepo: roleRepo,
	}
}

func (r *roleService) Create(c echo.Context, in *models.RoleWriteRequest) error {
	exist, err := r.roleRepo.IsEqualTypeExist(in.Type)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if exist {
		return helpers.Response(c, http.StatusBadRequest, "nama role sudah digunakan")
	}

	role := in.ToModelCreate()
	if err := r.roleRepo.Create(role); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil membuat role baru", role.ToResponse())
}

func (r *roleService) Update(c echo.Context, in *models.RoleWriteRequest) error {
	role, err := r.roleRepo.GetByIdWithSelectedFields(in.Id, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if role.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "role tidak ditemukan")
	}

	role, err = r.roleRepo.GetByTypeWithSelectedFields(in.Type, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if !role.IsEmpty() && role.Id != in.Id {
		return helpers.Response(c, http.StatusBadRequest, "nama role sudah digunakan")
	}

	role = in.ToModelUpdate()
	if err := r.roleRepo.Update(role); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil mengubah role", role.ToResponse())
}

func (r *roleService) GetListMaster(c echo.Context, search string) error {
	roles, err := r.roleRepo.GetAll(search)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	var res []*models.RoleResponse
	for _, r := range roles {
		res = append(res, r.ToResponse())
	}

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan semua role", res)
}

func (r *roleService) Delete(c echo.Context, id string) error {
	role, err := r.roleRepo.GetByIdWithSelectedFields(id, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if role.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "role tidak ditemukan")
	}

	if err := r.roleRepo.Delete(role); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil menghapus role")
}
