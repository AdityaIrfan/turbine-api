package services

import (
	"net/http"
	contract "turbine-api/contracts"
	"turbine-api/helpers"
	"turbine-api/models"

	"github.com/labstack/echo/v4"
)

type divisionService struct {
	divisionRepo contract.IDivisionRepository
}

func NewDivisionService(divisionRepo contract.IDivisionRepository) contract.IDivisionService {
	return &divisionService{
		divisionRepo: divisionRepo,
	}
}

func (d *divisionService) Create(c echo.Context, in *models.DivisionWriteRequest) error {
	exist, err := d.divisionRepo.IsEqualTypeExist(in.Type)
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if exist {
		return helpers.Response(c, http.StatusBadRequest, "division type already in use")
	}

	division := in.ToModelCreate()
	if err := d.divisionRepo.Create(division); err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	}

	return helpers.Response(c, http.StatusOK, "success create role", division.ToResponse())
}

func (d *divisionService) Update(c echo.Context, in *models.DivisionWriteRequest) error {
	division, err := d.divisionRepo.GetByIdWithSelectedFields(in.Id, "id")
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if division.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "division not found")
	}

	division, err = d.divisionRepo.GetByTypeWithSelectedFields(in.Type, "id")
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if !division.IsEmpty() && division.Id != in.Id {
		return helpers.Response(c, http.StatusBadRequest, "division type already in use")
	}

	division = in.ToModelUpdate()
	if err := d.divisionRepo.Update(division); err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	}

	return helpers.Response(c, http.StatusOK, "success update division", division.ToResponse())
}

func (d *divisionService) GetListMaster(c echo.Context, search string) error {
	divisions, err := d.divisionRepo.GetAll(search)
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	}

	var res []*models.DivisionResponse
	for _, d := range divisions {
		res = append(res, d.ToResponse())
	}

	return helpers.Response(c, http.StatusOK, "success get all divisions", res)
}

func (d *divisionService) Delete(c echo.Context, id string) error {
	division, err := d.divisionRepo.GetByIdWithSelectedFields(id, "id")
	if err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	} else if division.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "division not found")
	}

	if err := d.divisionRepo.Delete(division); err != nil {
		return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
	}

	return helpers.Response(c, http.StatusOK, "success delete division")
}
