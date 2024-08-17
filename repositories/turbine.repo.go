package repositories

import (
	"errors"
	"fmt"
	contract "pln/AdityaIrfan/turbine-api/contracts"
	"pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"
	"strings"

	"github.com/phuslu/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type turbineRepository struct {
	db *gorm.DB
}

func NewTurbineRepository(db *gorm.DB) contract.ITurbineRepository {
	return &turbineRepository{db: db}
}

func (t *turbineRepository) Create(turbine *models.Turbine) error {
	if err := t.db.Debug().
		Clauses(clause.Returning{}).
		Create(&turbine).
		Preload("PltaUnit.Plta").
		First(&turbine).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY TURBINE CREATING : " + err.Error()))
		return err
	}

	return nil
}
func (t *turbineRepository) GetByIdWithSelectedFields(id string, selectedFields string, preloads ...string) (*models.Turbine, error) {
	var turbine *models.Turbine

	db := t.db
	for _, p := range preloads {
		db = db.Preload(p)
	}

	if err := db.Debug().
		Select(selectedFields).
		Where("id = ?", id).
		First(&turbine).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		log.Error().Err(errors.New("ERROR QUERY TURBINE GETTING BY ID WITH SELECTED FIELDS : " + err.Error()))
		return nil, err
	}

	return turbine, nil
}
func (t *turbineRepository) GetAllWithPaginate(cursor *helpers.Cursor, selectedFields string) ([]*models.Turbine, *helpers.CursorPagination, error) {
	db := t.db
	alreadyWithPltaUnit := false

	if cursor.Search != "" {
		if strings.Contains(strings.ToLower(cursor.Search), "unit") {
			cursor.Search = strings.ReplaceAll(cursor.Search, "unit", "")
		} else if strings.Contains(strings.ToLower(cursor.Search), "unit ") {
			cursor.Search = strings.ReplaceAll(cursor.Search, "unit ", "")
		}

		db = db.Joins("LEFT JOIN plta_units ON plta_units.id = turbines.plta_unit_id").
			Joins("LEFT JOIN plta ON plta.id = plta_units.plta_id").
			Where("LOWER(plta_units.name) LIKE LOWER(?)", "%"+cursor.Search+"%").
			Or("LOWER(turbines.title) LIKE LOWER(?)", "%"+cursor.Search+"%").
			Or("LOWER(plta.name) LIKE LOWER(?)", "%"+cursor.Search+"%")
		alreadyWithPltaUnit = true
	}

	if cursor.StartDate != "" {
		db = db.Where("turbines.created_at >= ?", cursor.StartDate)
	}

	if cursor.EndDate != "" {
		db = db.Where("turbines.created_at <= ?", cursor.EndDate)
	}

	if cursor.Filter != "" {
		db = db.Where(fmt.Sprintf("%v = ?", cursor.Filter), cursor.FilterValue)
	}

	var total int64
	if err := db.Table("turbines").Count(&total).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY TURBINES TOTAL : " + err.Error())).Msg("")
		return nil, nil, err
	}

	var sortBy string
	switch strings.ToLower(cursor.SortBy) {
	case "title":
		sortBy = "turbines.title"
	case "pltaUnitName":
		if !alreadyWithPltaUnit {
			db = db.Joins("LEFT JOIN plta_units ON plta_units.id = turbines.plta_unit_id").
				Joins("LEFT JOIN plta ON plta_units.plta_id = plta.id")
		}

		selectedFields += ", plta.name || '- Unit ' || plta_units.name"
		sortBy = "plta_units.name"
	case "createdat":
		sortBy = "turbines.created_at"
	default:
		sortBy = "turbines.created_at"
	}

	var turbine []*models.Turbine
	if err := db.Debug().
		Select(selectedFields).
		Offset(cursor.CurrentPage*cursor.PerPage - cursor.PerPage).
		Limit(cursor.PerPage).
		Preload("PltaUnit.Plta").
		Order(fmt.Sprintf("%v %v", sortBy, cursor.SortOrder)).
		Find(&turbine).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY TURBINES LIST WITH PAGINATE : " + err.Error())).Msg("")
		return nil, nil, err
	}

	if total == 0 {
		return nil, &helpers.CursorPagination{}, nil
	}

	cursorPagination := cursor.GeneratePager(total)

	return turbine, cursorPagination, nil
}

func (t *turbineRepository) GetLatest() (*models.Turbine, error) {
	var turbine *models.Turbine

	if err := t.db.
		Preload("PltaUnit.Plta").
		Order("created_at desc").
		First(&turbine).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY LATEST TURBINE : " + err.Error())).Msg("")
		return nil, err
	}

	return turbine, nil
}

func (t *turbineRepository) Delete(turbine *models.Turbine) error {
	if err := t.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Table(turbine.TableName()).
			Where("id = ?", turbine.Id).
			Update("deleted_by", turbine.DeletedBy).Error; err != nil {
			log.Error().Err(errors.New("ERROR QUERY UPDATE deleted_by TURBINE: " + err.Error())).Msg("")
			return err
		}

		if err := tx.Where("id = ?", turbine.Id).Delete(&turbine).Error; err != nil {
			log.Error().Err(errors.New("ERROR QUERY DELETE TURBINE : " + err.Error())).Msg("")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (t *turbineRepository) GetTotal() (int64, error) {
	var counter int64

	if err := t.db.
		Table("turbines").
		Where("deleted_at IS NULL").
		Count(&counter).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY TOTAL USER : " + err.Error()))
		return 0, err
	}

	return counter, nil
}
