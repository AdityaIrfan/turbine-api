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
	userRepo     contract.IUserRepository
}

func NewDivisionService(divisionRepo contract.IDivisionRepository, userRepo contract.IUserRepository) contract.IDivisionService {
	return &divisionService{
		divisionRepo: divisionRepo,
		userRepo:     userRepo,
	}
}

func (d *divisionService) Create(c echo.Context, in *models.DivisionWriteRequest) error {
	exist, err := d.divisionRepo.IsEqualNameExist(in.Name)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if exist {
		return helpers.Response(c, http.StatusBadRequest, "division type already in use")
	}

	division := in.ToModelCreate()
	if err := d.divisionRepo.Create(division); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "success create role", division.ToResponse())
}

func (d *divisionService) Update(c echo.Context, in *models.DivisionWriteRequest) error {
	division, err := d.divisionRepo.GetByIdWithSelectedFields(in.Id, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if division.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "division not found")
	}

	divisionByType, err := d.divisionRepo.GetByNameWithSelectedFields(in.Name, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if !divisionByType.IsEmpty() && divisionByType.Id != in.Id {
		return helpers.Response(c, http.StatusBadRequest, "division type already in use")
	}

	division = in.ToModelUpdate()
	if err := d.divisionRepo.Update(division); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "success update division", division.ToResponse())
}

func (d *divisionService) GetListMaster(c echo.Context, search string) error {
	divisions, err := d.divisionRepo.GetAll(search)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	var res []*models.DivisionMasterResponse
	for _, d := range divisions {
		res = append(res, d.ToMasterResponse())
	}

	return helpers.Response(c, http.StatusOK, "success get all divisions", res)
}

func (d *divisionService) GetListWithPaginate(c echo.Context, cursor *helpers.Cursor) error {
	divisions, pagination, err := d.divisionRepo.GetAllWithPaginate(cursor)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	var divisionRes []*models.DivisionResponse
	for _, division := range divisions {
		divisionRes = append(divisionRes, division.ToResponse())
	}

	return helpers.Response(c, http.StatusOK, "success get division list", divisionRes, pagination)
}

func (d *divisionService) Delete(c echo.Context, in *models.DivisionWriteRequest) error {
	division, err := d.divisionRepo.GetByIdWithSelectedFields(in.Id, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if division.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "division not found")
	}

	if err := d.divisionRepo.Delete(division); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "success delete division")
}
