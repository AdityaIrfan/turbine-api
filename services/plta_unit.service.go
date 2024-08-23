package services

import (
	"fmt"
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
	userRepo     contract.IUserRepository
}

func NewPltaUnitService(
	pltaUnitRepo contract.IPltaUnitRepository,
	pltaRepo contract.IPltaRepository,
	userRepo contract.IUserRepository) contract.IPltaUnitService {
	return &pltaUnitService{
		pltaUnitRepo: pltaUnitRepo,
		pltaRepo:     pltaRepo,
		userRepo:     userRepo,
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

func (p *pltaUnitService) GetListMaster(c echo.Context, in *models.PltaGetListMasterRequest) error {
	user, err := p.userRepo.GetByIdWithSelectedFields(in.UserId, "id, radius_status")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		return helpers.ResponseForbiddenAccess(c)
	}

	units, err := p.pltaUnitRepo.GetAll(in.Search)
	if err != nil {
		return helpers.ResponseForbiddenAccess(c)
	}

	var pltaResponse []*models.PltaResponseMaster
	var pltaMap = make(map[string]*models.Plta)
	for _, unit := range units {
		var plta *models.Plta
		if value, ok := pltaMap[unit.PltaId]; ok {
			plta = value
		} else {
			res, err := p.pltaRepo.GetByIdWithSelectedFields(unit.PltaId, "name, lat, long, radius_status, radius, radius_type")
			if err != nil {
				return helpers.ResponseUnprocessableEntity(c)
			}
			pltaMap[unit.PltaId] = res
			plta = res
		}

		if !user.RadiusStatus {
			plta.RadiusStatus = false
		}

		pltaResponse = append(pltaResponse, &models.PltaResponseMaster{
			Id:           unit.Id,
			Name:         fmt.Sprintf("%v - Unit %v", plta.Name, unit.Name),
			Lat:          plta.Lat,
			Long:         plta.Long,
			RadiusStatus: plta.RadiusStatus,
			Radius:       plta.Radius,
			RadiusType:   plta.RadiusType,
		})
	}

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan semua plta master", pltaResponse)
}
