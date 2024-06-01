package repositories

import (
	"errors"
	"strings"
	contract "turbine-api/contracts"
	"turbine-api/models"

	"github.com/phuslu/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type divisionRepository struct {
	db *gorm.DB
}

func NewDivisionRepository(db *gorm.DB) contract.IDivisionRepository {
	return &divisionRepository{
		db: db,
	}
}

func (d *divisionRepository) Create(division *models.Division) error {
	if err := d.db.Create(&division).Clauses(clause.Returning{}).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY DIVISION INSERT : " + err.Error()))
		return err
	}

	return nil
}

func (d *divisionRepository) Update(division *models.Division) error {
	if err := d.db.Updates(&division).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY DIVISION UPDATE : " + err.Error()))
		return err
	}

	return nil
}

func (d *divisionRepository) GetAll(search string) ([]*models.Division, error) {
	var divisions []*models.Division

	db := d.db
	if search != "" {
		db = db.Where("LOWER(type) LIKES ?", "'%"+strings.ToLower(search)+"%'")
	}

	if err := db.Find(&divisions).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY ALL DIVISIONS : " + err.Error()))
		return nil, err
	}

	return divisions, nil
}

func (d *divisionRepository) GetById(id string) (*models.Division, error) {
	var division *models.Division

	if err := d.db.Where("id = ?", id).First(&division).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY ROLE BY ID : " + err.Error()))
		return nil, err
	}

	return division, nil
}

func (d *divisionRepository) GetByIdWithSelectedFields(id string, selectedFields string) (*models.Division, error) {
	var division *models.Division

	if err := d.db.Select(selectedFields).Where("id = ?", id).First(&division).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY DIVISION BY ID WITH SELECTED FIELDS : " + err.Error()))
		return nil, err
	}

	return division, nil
}

func (d *divisionRepository) GetByTypeWithSelectedFields(divisionTYpe models.DivisionType, selectedFields string) (*models.Division, error) {
	var division *models.Division

	if err := d.db.Select(selectedFields).Where("type = ?", divisionTYpe).First(&division).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY DIVISION TYPE ID WITH SELECTED FIELDS : " + err.Error()))
		return nil, err
	}

	return division, nil
}

func (d *divisionRepository) IsEqualTypeExist(divisionType models.DivisionType) (bool, error) {
	var division *models.Division

	err := d.db.Where("LOWER(type) LIKES ?", "'%"+strings.ToLower(string(divisionType))+"%'").First(&division).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		log.Error().Err(errors.New("ERROR QUERY DIVISION IsEqualTypeExist : " + err.Error()))
		return false, err
	}

	return true, nil
}

func (d *divisionRepository) Delete(division *models.Division) error {
	if err := d.db.Delete(&division).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY DIVISION DELETE : " + err.Error()))
		return err
	}

	return nil
}
