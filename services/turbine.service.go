package services

import (
	"net/http"

	contract "pln/AdityaIrfan/turbine-api/contracts"
	"pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"

	"github.com/labstack/echo/v4"
)

type turbineService struct {
	turbineRepo contract.ITurbineRepository
	towerRepo   contract.ITowerRepository
}

func NewTurbineService(
	turbineRepo contract.ITurbineRepository,
	towerRepo contract.ITowerRepository) contract.ITurbineService {
	return &turbineService{
		turbineRepo: turbineRepo,
		towerRepo:   towerRepo,
	}
}

func (t *turbineService) Create(c echo.Context, in *models.TurbineWriteRequest) error {
	turbine := in.ToModelCreate()
	tower, err := t.towerRepo.GetByIdWithSelectedFields(turbine.TowerId, "id, name, unit_number")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if tower.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "tower tidak ditemukan")
	}

	turbine.Tower = tower

	if err := t.turbineRepo.Create(turbine); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil menambahkan data turbine baru", turbine.ToResponse())
}

func (t *turbineService) GetDetail(c echo.Context, id string) error {
	turbine, err := t.turbineRepo.GetByIdWithSelectedFields(id, "*", "Tower")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if turbine.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "data turbine tidak ditemukan")
	}

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan data turbine", turbine.ToResponse())
}

func (t *turbineService) GetListWithPaginate(c echo.Context, cursor *helpers.Cursor) error {
	turbines, pagination, err := t.turbineRepo.GetAllWithPaginate(cursor, "")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	var turbineResponse = []*models.TurbineResponseList{}
	for _, turbine := range turbines {
		turbineResponse = append(turbineResponse, turbine.ToResponseList())
	}

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan semua data turbine", turbineResponse, pagination)
}

func (t *turbineService) GetLatest(c echo.Context) error {
	turbine, err := t.turbineRepo.GetLatest()
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}
	if turbine.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "data tidak ditemukan")
	}

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan data turbine terakhir", turbine.ToResponse())
}
