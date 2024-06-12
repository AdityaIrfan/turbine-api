package repositories

import (
	"errors"
	"fmt"
	"strings"
	contract "turbine-api/contracts"
	"turbine-api/helpers"
	"turbine-api/models"

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
	if err := t.db.
		Clauses(clause.Returning{}).
		Create(&turbine).
		Preload("User", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id, name")
		}).First(&turbine).Error; err != nil {
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

	if err := db.
		Select(selectedFields).
		Where("id = ?", id).
		Preload("User", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id, name")
		}).
		First(&turbine).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY TURBINE GETTING BY ID WITH SELECTED FIELDS : " + err.Error()))
		return nil, err
	}

	return turbine, nil
}
func (t *turbineRepository) GetAllWithPaginate(cursor *helpers.Cursor, selectedFields string) ([]*models.Turbine, *helpers.CursorPagination, error) {
	db := t.db
	alreadyWithTower := false

	if cursor.Search != "" {
		db = db.Joins("LEFT JOIN towers ON towers.id = turbines.tower_id").
			Where("LOWER(towers.name) LIKE LOWER(?)", "%"+cursor.Search+"%").
			Or("LOWER(towers.unit_number) LIKE LOWER(?)", "%"+cursor.Search+"%")
		alreadyWithTower = true
	}

	var total int64
	if err := db.Table("turbines").Count(&total).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY TURBINES TOTAL : " + err.Error())).Msg("")
		return nil, nil, err
	}

	var sortBy string
	switch strings.ToLower(cursor.SortBy) {
	case "towername":
		if !alreadyWithTower {
			db = db.Joins("LEFT JOIN towers ON towers.id = turbines.tower_id")
		}

		sortBy = "towers.name"
	case "createdat":
		sortBy = "created_at"
	default:
		sortBy = "created_at"
	}

	var turbine []*models.Turbine
	if err := db.Debug().
		Offset(cursor.CurrentPage - 1).
		Limit(cursor.PerPage).
		Preload("Tower").
		Order(fmt.Sprintf("%v %v", sortBy, cursor.SortOrder)).
		Find(&turbine).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY TURBINES LIST WITH PAGINATE : " + err.Error())).Msg("")
		return nil, nil, err
	}

	cursorPagination := cursor.GeneratePager(total)

	return turbine, cursorPagination, nil
}
