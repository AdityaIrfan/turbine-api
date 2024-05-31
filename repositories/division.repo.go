package repositories

import (
	"errors"
	"strings"
	contract "turbine-api/contracts"
	"turbine-api/models"

	"gorm.io/gorm"
)

type divisionRepository struct {
	db *gorm.DB
}

func NewDivisionRepository(db *gorm.DB) contract.IDivisionRepository {
	return &divisionRepository{
		db: db,
	}
}

func (r *divisionRepository) GetAll(search string) ([]*models.Division, error) {
	var divisions []*models.Division

	db := r.db
	if search != "" {
		db = db.Where("LOWER(type) LIKES ?", "'%"+strings.ToLower(search)+"%'")
	}

	err := db.Find(&divisions).Error
	return divisions, err
}

func (r *divisionRepository) GetById(id string) (*models.Division, error) {
	var division *models.Division

	db := r.db
	if err := db.Where("id = ?", id).First(&division).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return division, nil
}

func (r *divisionRepository) IsEqualTypeExist(divisionType models.DivisionType) (bool, error) {
	var division *models.Division

	db := r.db

	err := db.Where("LOWER(type) LIKES ?", "'%"+strings.ToLower(string(divisionType))+"%'").First(&division).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
