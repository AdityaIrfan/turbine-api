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
	towerByName, err := t.towerRepo.GetByEqualNameWithSelectedFields(in.Name, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if !towerByName.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "name already in use")
	}

	tower := in.ToModelCreate()
	if err := t.towerRepo.Create(tower); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "success create tower", tower.ToResponse())
}

func (t *towerService) Update(c echo.Context, in *models.TowerWriteRequest) error {
	tower, err := t.towerRepo.GetByIdWithSelectedFields(in.Id, "*")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if tower.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "tower not found")
	}

	if in.Name != tower.Name {
		towerByName, err := t.towerRepo.GetByEqualNameWithSelectedFields(in.Name, "id")
		if err != nil {
			return helpers.ResponseUnprocessableEntity(c)
		} else if !towerByName.IsEmpty() {
			return helpers.Response(c, http.StatusBadRequest, "name already in use")
		}

		tower := in.ToModelCreate()
		if err := t.towerRepo.Update(tower); err != nil {
			return helpers.ResponseUnprocessableEntity(c)
		}
	}

	return helpers.Response(c, http.StatusOK, "success create tower", tower.ToResponse())
}

func (t *towerService) GetListMaster(c echo.Context, search string) error {
	towers, err := t.towerRepo.GetAll(search)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	var towerResponse []*models.TowerResponse
	for _, tower := range towers {
		towerResponse = append(towerResponse, tower.ToResponse())
	}

	return helpers.Response(c, http.StatusOK, "success get list master tower", towerResponse)
}

func (t *towerService) Delete(c echo.Context, in *models.TowerWriteRequest) error {
	tower, err := t.towerRepo.GetByIdWithSelectedFields(in.Id, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if tower.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "tower not found")
	}

	if err := t.towerRepo.Delete(tower); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "success delete tower")
}
