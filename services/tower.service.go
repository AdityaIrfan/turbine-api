package services

import (
	"net/http"
	contract "turbine-api/contracts"
	"turbine-api/helpers"
	"turbine-api/models"

	"github.com/labstack/echo/v4"
)

type towerService struct {
	towerRepo contract.ITowerRepository
}

func NewTowerService(towerRepo contract.ITowerRepository) contract.ITowerService {
	return &towerService{
		towerRepo: towerRepo,
	}
}

func (t *towerService) Create(c echo.Context, in *models.TowerWriteRequest) error {
	towerByName, err := t.towerRepo.GetByEqualNameAndUnitNumberWithSelectedFields(in.Name, in.UnitNumber, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if !towerByName.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "nama tower sudah digunakan")
	}

	tower := in.ToModelCreate()
	if err := t.towerRepo.Create(tower); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil membuat tower baru", tower.ToResponse())
}

func (t *towerService) Update(c echo.Context, in *models.TowerWriteRequest) error {
	tower, err := t.towerRepo.GetByIdWithSelectedFields(in.Id, "*")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if tower.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "tower tidak ditemukan")
	}

	if in.Name != tower.Name {
		towerByName, err := t.towerRepo.GetByEqualNameAndUnitNumberWithSelectedFields(in.Name, in.UnitNumber, "id")
		if err != nil {
			return helpers.ResponseUnprocessableEntity(c)
		} else if !towerByName.IsEmpty() {
			return helpers.Response(c, http.StatusBadRequest, "nama tower sudah digunakan")
		}

		tower := in.ToModelCreate()
		if err := t.towerRepo.Update(tower); err != nil {
			return helpers.ResponseUnprocessableEntity(c)
		}
	}
	if in.UnitNumber != tower.UnitNumber {
		tower.UnitNumber = in.UnitNumber
	}

	if err := t.towerRepo.Update(tower); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil mengubah tower", tower.ToResponse())
}

func (t *towerService) GetListMaster(c echo.Context, search string) error {
	towers, err := t.towerRepo.GetAll(search)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	var towerResponse []*models.TowerResponseMaster
	for _, tower := range towers {
		towerResponse = append(towerResponse, tower.ToResponseMaster())
	}

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan semua tower master", towerResponse)
}

func (t *towerService) Delete(c echo.Context, id string) error {
	tower, err := t.towerRepo.GetByIdWithSelectedFields(id, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if tower.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "tower tidak ditemukan")
	}

	if err := t.towerRepo.Delete(tower); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil menghapus tower")
}
