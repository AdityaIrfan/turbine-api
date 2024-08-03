package services

import (
	"net/http"
	contract "pln/AdityaIrfan/turbine-api/contracts"
	"pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"
	"strings"

	"github.com/labstack/echo/v4"
)

type pltaUnitService struct {
	pltaUnitRepo contract.IPltaUnitRepository
	pltaRepo     contract.IPltaRepository
}

func NewPltaUnitService(pltaUnitRepo contract.IPltaUnitRepository, pltaRepo contract.IPltaRepository) contract.IPltaUnitService {
	return &pltaUnitService{
		pltaUnitRepo: pltaUnitRepo,
		pltaRepo:     pltaRepo,
	}
}

func (p *pltaUnitService) CreateOrUpdate(c echo.Context, in *models.PltaUnitCreateOrUpdate) error {
	plta, err := p.pltaRepo.GetByIdWithSelectedFields(in.PltaId, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if plta.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "plta tidak ditemukan")
	}

	pltaUnits := in.ToModelCreateOrUpdate()
	pltaUnits, err = p.pltaUnitRepo.CreateOrUpdate(pltaUnits)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			return helpers.Response(c, http.StatusNotFound, err.Error())
		} else if strings.Contains(err.Error(), "duplikat") {
			return helpers.Response(c, http.StatusBadRequest, err.Error())
		}
		return helpers.ResponseUnprocessableEntity(c)
	}

	var res = []*models.PltaUnitResponse{}
	for _, unit := range pltaUnits {
		res = append(res, unit.ToResponse())
	}

	return helpers.Response(c, http.StatusOK, "berhasil menambah atau mengubah plta unit", res)
}

func (p *pltaUnitService) Delete(c echo.Context, in *models.PltaUnitWriteRequest) error {
	pltaUnit, err := p.pltaUnitRepo.GetByIdAndSelectedFields(in.Id, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if pltaUnit.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "plta unit tidak ditemukan")
	}

	pltaUnit.DeletedBy = in.WrittenBy
	if err := p.pltaUnitRepo.Delete(pltaUnit); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil menghapus plta unit")
}
