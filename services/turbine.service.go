package services

import (
	"net/http"
	contract "turbine-api/contracts"
	"turbine-api/helpers"
	"turbine-api/models"

	"github.com/labstack/echo/v4"
)

type turbineService struct {
	turbineRepo contract.ITurbineRepository
}

func NewTurbineService(turbineRepo contract.ITurbineRepository) contract.ITurbineService {
	return &turbineService{
		turbineRepo: turbineRepo,
	}
}

func (t *turbineService) Create(c echo.Context, in *models.TurbineWriteRequest) error {
	turbine := in.ToModelCreate()
	// if err := t.turbineRepo.Create(turbine); err != nil {
	// return helpers.ResponseUnprocessableEntity(c)
	// }

	return helpers.Response(c, http.StatusOK, "berhasil menambahkan data turbine baru", turbine.ToResponse())
}

func (t *turbineService) GetDetail(c echo.Context, id string) error {
	turbine, err := t.turbineRepo.GetById(id, "*")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if turbine.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "data turbine tidak ditemukan")
	}

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan data turbine")
}

func (t *turbineService) GetListWithPaginate(c echo.Context, cursor *helpers.Cursor) error {
	turbines, pagination, err := t.turbineRepo.GetAllWithPaginate(cursor, "")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan semua data turbine", turbines, pagination)
}
