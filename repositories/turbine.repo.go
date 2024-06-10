package repositories

import (
	contract "turbine-api/contracts"
	"turbine-api/helpers"
	"turbine-api/models"

	"gorm.io/gorm"
)

type turbineRepository struct {
	db *gorm.DB
}

func NewTurbineRepository(db *gorm.DB) contract.ITurbineRepository {
	return &turbineRepository{db: db}
}

func (t *turbineRepository) Create(turbine *models.Turbine) error {
	return nil
}
func (t *turbineRepository) GetById(id string, selectedFields string) (*models.Turbine, error) {
	return nil, nil
}
func (t *turbineRepository) GetAllWithPaginate(cursor *helpers.Cursor, selectedFields string) ([]*models.Turbine, *helpers.CursorPagination, error) {
	return nil, nil, nil
}
