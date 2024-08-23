package repositories

import (
	"errors"
	"fmt"
	contract "pln/AdityaIrfan/turbine-api/contracts"
	"pln/AdityaIrfan/turbine-api/models"
	"strings"

	"github.com/oklog/ulid/v2"
	"github.com/phuslu/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type pltaUnitRepo struct {
	db *gorm.DB
}

func NewPltaUnitRepo(db *gorm.DB) contract.IPltaUnitRepository {
	return &pltaUnitRepo{
		db: db,
	}
}

func (p *pltaUnitRepo) GetByIdAndSelectedFields(id, selectedFields string) (*models.PltaUnit, error) {
	var pltaUnit *models.PltaUnit

	if err := p.db.Where("id = ?", id).First(&pltaUnit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY GETTING PLTA UNIT BY ID AND SELECTED FIELDS : " + err.Error())).Msg("")
		return nil, err
	}

	return pltaUnit, nil
}

func (p *pltaUnitRepo) CreateOrUpdate(pltaUnits []*models.PltaUnit) ([]*models.PltaUnit, error) {
	var pltaId string

	if err := p.db.Transaction(func(tx *gorm.DB) error {
		for _, unit := range pltaUnits {
			if unit.Id == "" {
				var check *models.PltaUnit
				if err := tx.Where("name = ?", unit.Name).Select("id").First(&check).Error; err != nil {
					if !errors.Is(err, gorm.ErrRecordNotFound) {
						log.Error().Err(errors.New("ERROR QUERY GETTING PLTA UNIT BY ID ON CREATING OR UPDATING PLTA UNIT : " + err.Error())).Msg("")
						return err
					}
				} else if check != nil {
					return fmt.Errorf("duplikat nama plta unit, nama %s sudah digunakan", unit.Name)
				}

				unit.Id = ulid.Make().String()
				if err := tx.
					Clauses(clause.Returning{}).
					Create(&unit).Error; err != nil {
					log.Error().Err(errors.New("ERROR QUERY CREATING PLTA UNIT : " + err.Error())).Msg("")
					return err
				}

				if pltaId == "" {
					pltaId = unit.PltaId
				}
			} else {
				var check *models.PltaUnit
				if err := tx.Where("id = ?", unit.Id).Select("id").First(&check).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return fmt.Errorf("plta unit dengan id %s tidak ditemukan", unit.Id)
					}
					log.Error().Err(errors.New("ERROR QUERY GETTING PLTA UNIT BY ID ON CREATING OR UPDATING PLTA UNIT : " + err.Error())).Msg("")
					return err
				}

				if !unit.Status {
					if err := tx.
						Table("plta_units").
						Where("id = ?", unit.Id).
						Update("status", false).Error; err != nil {
						log.Error().Err(errors.New("ERROR QUERY UPDATING STATUS ON UPDATING PLTA UNIT : " + err.Error())).Msg("")
						return err
					}
				}

				if err := tx.
					Clauses(clause.Returning{}).
					Updates(&unit).Error; err != nil {
					log.Error().Err(errors.New("ERROR QUERY UPDATING PLTA UNITS : " + err.Error())).Msg("")
					return err
				}

				if pltaId == "" {
					pltaId = unit.PltaId
				}
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	if err := p.db.Where("plta_id = ?", pltaId).Clauses(clause.Returning{}).Find(&pltaUnits).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY GETTING ALL PLTA UNITS BY PLTA ID ON CREATING AND UPDATING PLTA UNITS : " + err.Error())).Msg("")
		return nil, err
	}

	return pltaUnits, nil
}

func (p *pltaUnitRepo) Delete(pltaUnit *models.PltaUnit) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Table("plta_units").
			Where("id = ?", pltaUnit.Id).
			Updates(map[string]interface{}{
				"status":     false,
				"deleted_by": pltaUnit.DeletedBy,
			}).Error; err != nil {
			log.Error().Err(errors.New("ERROR QUERY UPDATING DELETED BY ON DELETING PLTA UNITS : " + err.Error())).Msg("")
			return err
		}

		if err := tx.Delete(&pltaUnit).Error; err != nil {
			log.Error().Err(errors.New("ERROR QUERY DELETING PLTA UNIT : " + err.Error())).Msg("")
			return err
		}

		return nil
	})
}

func (p *pltaUnitRepo) GetByIdWithPreloads(id string, preloads ...string) (*models.PltaUnit, error) {
	db := p.db.Debug()

	for _, p := range preloads {
		db = db.Preload(p)
	}

	var pltaUnit *models.PltaUnit
	if err := db.Where("id = ?", id).First(&pltaUnit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY GETTING PLTA UNITS BY ID : " + err.Error())).Msg("")
		return nil, err
	}

	return pltaUnit, nil
}

func (p *pltaUnitRepo) GetAll(search string) ([]*models.PltaUnit, error) {
	db := p.db

	if strings.Contains(strings.ToLower(search), "unit") {
		search = strings.ReplaceAll(strings.ToLower(search), "unit", "")
		search = strings.Trim(search, " ")
	}

	var units []*models.PltaUnit

	if err := db.
		Joins("LEFT JOIN plta on plta.id = plta_units.plta_id").
		Where("plta_units.status", true).
		Where("plta.status", true).
		Where(
			p.db.Where("LOWER(plta_units.name) LIKE LOWER(?)", "%"+search+"%").
				Or("LOWER(plta.name) LIKE LOWER(?)", "%"+search+"%"),
		).
		Find(&units).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY GET ALL PLTA UNITS : " + err.Error())).Msg("")
		return nil, err
	}

	return units, nil
}
