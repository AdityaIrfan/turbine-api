package services

import (
	"fmt"
	"net/http"
	"strconv"

	contract "pln/AdityaIrfan/turbine-api/contracts"
	"pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"

	"github.com/labstack/echo/v4"
)

type turbineService struct {
	turbineRepo  contract.ITurbineRepository
	pltaUnitRepo contract.IPltaUnitRepository
	userRepo     contract.IUserRepository
}

func NewTurbineService(
	turbineRepo contract.ITurbineRepository,
	pltaUnitRepo contract.IPltaUnitRepository,
	userRepo contract.IUserRepository) contract.ITurbineService {
	return &turbineService{
		turbineRepo:  turbineRepo,
		pltaUnitRepo: pltaUnitRepo,
		userRepo:     userRepo,
	}
}

func (t *turbineService) Create(c echo.Context, in *models.TurbineWriteRequest) error {
	user, err := t.userRepo.GetByIdWithSelectedFields(in.WrittenBy, "id, name, radius_status")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		return helpers.ResponseForbiddenAccess(c)
	}

	pltaUnit, err := t.pltaUnitRepo.GetByIdWithPreloads(in.TowerId, "Plta")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if pltaUnit.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "plta unit tidak ditemukan")
	} else if !pltaUnit.IsActive() {
		return helpers.Response(c, http.StatusForbidden, "plta unit tidak aktif, tidak dapat memasukkan data")
	} else if !pltaUnit.Plta.IsActive() {
		return helpers.Response(c, http.StatusForbidden, "plta tidak aktif, tidak dapat memasukkan data")
	}

	// if user.IsRadiusStatusActive() {
	// 	if pltaUnit.Plta.IsRadiusStatusActive() {
	// 		withinArea, err := helpers.IsIPWithinRadius(
	// 			c.Get("IP").(string),
	// 			pltaUnit.Plta.Lat,
	// 			pltaUnit.Plta.Long,
	// 			pltaUnit.Plta.GetRadiusInKilometer())
	// 		if err != nil {
	// 			return helpers.ResponseUnprocessableEntity(c)
	// 		} else if !withinArea {
	// 			return helpers.Response(c, http.StatusForbidden, "diluar jangkauan")
	// 		}
	// 	}
	// }

	turbine := in.ToModelCreate()
	if err := t.turbineRepo.Create(turbine); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	turbine.CreatedBy = user.Name

	return helpers.Response(c, http.StatusOK, "berhasil menambahkan data turbine baru", turbine.ToResponse())
}

func (t *turbineService) GetDetail(c echo.Context, id string) error {
	turbine, err := t.turbineRepo.GetByIdWithSelectedFields(id, "*", "PltaUnit.Plta")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if turbine.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "data turbine tidak ditemukan")
	}

	// created by
	user, err := t.userRepo.GetByIdWithSelectedFields(turbine.CreatedBy, "name")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		user = &models.User{}
	}

	turbine.CreatedBy = user.Name

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan data turbine", turbine.ToResponse())
}

func (t *turbineService) GetListWithPaginate(c echo.Context, cursor *helpers.Cursor) error {
	turbines, pagination, err := t.turbineRepo.GetAllWithPaginate(cursor, "turbines.id, turbines.title, turbines.created_at, turbines.plta_unit_id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	var turbineResponse = []*models.TurbineResponseList{}
	var pltaUnitsMap = make(map[string]*models.PltaUnit)
	for _, turbine := range turbines {
		if pltaUnit, ok := pltaUnitsMap[turbine.PltaUnitId]; ok {
			turbine.PltaUnit = pltaUnit
		} else {
			pltaUnit, _ := t.pltaUnitRepo.GetByIdWithPreloads(turbine.PltaUnitId, "Plta")
			turbine.PltaUnit = pltaUnit
			pltaUnitsMap[pltaUnit.Id] = pltaUnit
		}

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

	// created by
	user, err := t.userRepo.GetByIdWithSelectedFields(turbine.CreatedBy, "name")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		user = &models.User{}
	}

	turbine.CreatedBy = user.Name

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan data turbine terakhir", turbine.ToResponse())
}

func (t *turbineService) Delete(c echo.Context, in *models.TurbineWriteRequest) error {
	turbine, err := t.turbineRepo.GetByIdWithSelectedFields(in.Id, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if turbine.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "data turbine tidak ditemukan")
	}

	turbine.DeletedBy = in.WrittenBy
	if err := t.turbineRepo.Delete(turbine); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil menghapus data turbine", nil)
}

func (t *turbineService) DownloadReport(c echo.Context, id string) error {
	turbine, err := t.turbineRepo.GetByIdWithSelectedFields(id, "*", "PltaUnit.Plta")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if turbine.IsEmpty() {
		return helpers.Response(c, http.StatusNotFound, "data turbine tidak ditemukan")
	}

	// created by
	user, err := t.userRepo.GetByIdWithSelectedFields(turbine.CreatedBy, "name")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		user = &models.User{}
	}

	turbine.CreatedBy = user.Name

	reportBytes, err := turbine.GenerateReport()
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	// Set headers
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s [%s].pdf", turbine.Title, turbine.CreatedAt.Format("2006-01-02")))
	c.Response().Header().Set("Content-Length", strconv.Itoa(len(reportBytes)))

	return c.Blob(http.StatusOK, "application/pdf", reportBytes)
}
