package repositories

import (
	"errors"
	contract "pln/AdityaIrfan/turbine-api/contracts"
	"pln/AdityaIrfan/turbine-api/models"

	"github.com/phuslu/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type towerRepository struct {
	db *gorm.DB
}

func NewTowerRepository(db *gorm.DB) contract.ITowerRepository {
	return &towerRepository{
		db: db,
	}
}

func (t *towerRepository) Create(tower *models.Tower) error {
	if err := t.db.
		Clauses(clause.Returning{}).
		Create(&tower).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY TOWER CREATING : " + err.Error())).Msg("")
		return err
	}

	return nil
}

func (t *towerRepository) Update(tower *models.Tower) error {
	if err := t.db.
		Updates(&tower).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY TOWER UPDATING : " + err.Error())).Msg("")
		return err
	}

	return nil
}

func (t *towerRepository) GetByIdWithSelectedFields(id string, selectedFields string) (*models.Tower, error) {
	var tower *models.Tower

	if err := t.db.
		Select(selectedFields).
		Where("id = ?", id).
		First(&tower).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY GET TOWER BY ID WITH SELECTED FIELDS : " + err.Error())).Msg("")
		return nil, err
	}

	return tower, nil
}

func (t *towerRepository) GetByEqualNameAndUnitNumberWithSelectedFields(name, unitNumber string, selectedFields string) (*models.Tower, error) {
	var tower *models.Tower

	if err := t.db.
		Select(selectedFields).
		Where("LOWER(name) = LOWER(?)", name).
		Where("LOWER(unit_number) = LOWER(?)", unitNumber).
		First(&tower).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY GET TOWER BY ID WITH SELECTED FIELDS : " + err.Error())).Msg("")
		return nil, err
	}

	return tower, nil
}

func (t *towerRepository) GetAll(search string) ([]*models.Tower, error) {
	var towers []*models.Tower

	if err := t.db.
		Where("LOWER(name) LIKE LOWER(?)", "%"+search+"%").
		Or("LOWER(unit_number) LIKE LOWER(?)", "%"+search+"%").
		Find(&towers).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY GET ALL TOWERS : " + err.Error())).Msg("")
		return nil, err
	}

	return towers, nil
}

func (t *towerRepository) Delete(tower *models.Tower) error {
	if err := t.db.Delete(&tower).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY TOWER DELETING : " + err.Error())).Msg("")
		return err
	}

	return nil
}
