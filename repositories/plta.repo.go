package repositories

import (
	"errors"
	"fmt"
	contract "pln/AdityaIrfan/turbine-api/contracts"
	"pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"
	"strings"
	"time"

	"github.com/phuslu/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type pltaRepository struct {
	db *gorm.DB
}

func NewPltaRepository(db *gorm.DB) contract.IPltaRepository {
	return &pltaRepository{
		db: db,
	}
}

func (t *pltaRepository) Create(plta *models.Plta) error {
	if err := t.db.Debug().
		Clauses(clause.Returning{}).
		Create(&plta).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY PLTA CREATING : " + err.Error())).Msg("")
		return err
	}

	return nil
}

func (t *pltaRepository) Update(plta *models.Plta) error {
	return t.db.Debug().Transaction(func(tx *gorm.DB) error {
		if !plta.Status {
			if err := tx.Table("plta").Where("id = ?", plta.Id).
				Update("status", false).Error; err != nil {
				log.Error().Err(errors.New("ERROR QUERY UPDATING PLTA STATUS TO BE FALSE : " + err.Error())).Msg("")
				return err
			}

			if err := tx.
				Table("plta_units").
				Where("plta_id = ?", plta.Id).
				Where("status = true").
				Update("status", false).Error; err != nil {
				log.Error().Err(errors.New("ERROR QUERY PLTA UNITS UPDATED STATUS TO FALSE : " + err.Error())).Msg("")
				return err
			}
		}

		if !plta.RadiusStatus {
			if err := tx.
				Table("plta").
				Where("id = ?", plta.Id).
				Update("radius_status", false).Error; err != nil {
				log.Error().Err(errors.New("ERROR QUERY PLTA RADIUS STATUS TO BE FALSE : " + err.Error())).Msg("")
				return err
			}
		}

		if err := tx.
			Clauses(clause.Returning{}).
			Updates(&plta).
			Preload("PltaUnits").
			First(&plta).Error; err != nil {
			log.Error().Err(errors.New("ERROR QUERY PLTA UPDATING : " + err.Error())).Msg("")
			return err
		}
		return nil
	})
}

func (t *pltaRepository) GetByIdWithSelectedFields(id string, selectedFields string) (*models.Plta, error) {
	var plta *models.Plta

	if err := t.db.
		Select(selectedFields).
		Where("id = ?", id).
		First(&plta).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY GET PLTA BY ID WITH SELECTED FIELDS : " + err.Error())).Msg("")
		return nil, err
	}

	return plta, nil
}

func (t *pltaRepository) GetByEqualNameWithSelectedFields(name, selectedFields string) (*models.Plta, error) {
	var plta *models.Plta

	if err := t.db.
		Select(selectedFields).
		Where("LOWER(name) = LOWER(?)", name).
		First(&plta).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY GET PLTA BY ID WITH SELECTED FIELDS : " + err.Error())).Msg("")
		return nil, err
	}

	return plta, nil
}

func (t *pltaRepository) GetAll(search string) ([]*models.Plta, error) {
	var pltas []*models.Plta

	if strings.Contains(strings.ToLower(search), "unit") {
		search = strings.ReplaceAll(strings.ToLower(search), "unit", "")
		search = strings.Trim(search, " ")
	}

	if err := t.db.Select("plta.*").Debug().
		Joins("LEFT JOIN plta_units ON plta_units.plta_id = plta.id").
		Where(
			t.db.Where("LOWER(plta.name) LIKE LOWER(?)", "%"+search+"%").
				Or("LOWER(plta_units.name) LIKE LOWER(?)", "%"+search+"%"),
		).
		Where("plta.status = true").
		Preload("PltaUnits", func(db *gorm.DB) *gorm.DB {
			db = db.
				Where("status = true").
				Select("id, plta_id, name")

			if search != "" {
				db = db.Where("LOWER(name) LIKE ?", "%"+search+"%")
			}

			return db
		}).
		Find(&pltas).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY GET ALL PLTA : " + err.Error())).Msg("")
		return nil, err
	}

	return pltas, nil
}

func (t *pltaRepository) Delete(plta *models.Plta) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		if plta.Status {
			if err := tx.Table("plta").Where("id = ?", plta.Id).
				Update("status", false).Error; err != nil {
				log.Error().Err(errors.New("ERROR QUERY UPDATING PLTA STATUS TO BE FALSE : " + err.Error())).Msg("")
				return err
			}
		}

		if err := tx.
			Table("plta").
			Where("id = ?", plta.Id).
			Update("deleted_by", plta.DeletedBy).Error; err != nil {
			log.Error().Err(errors.New("ERROR QUERY UPDATING DELETED BY ON DELETING PLTA : " + err.Error())).Msg("")
			return err
		}

		if err := tx.
			Table("plta_units").
			Where("plta_id = ?", plta.Id).
			Updates(map[string]interface{}{
				"status":     false,
				"deleted_at": time.Now(),
				"deleted_by": plta.DeletedBy,
			}).Error; err != nil {
			log.Error().Err(errors.New("ERROR QUERY UPDATING DELETED BY PLTA UNITS ON DELETING PLTA : " + err.Error())).Msg("")
			return err
		}

		if err := tx.Delete(&plta).Error; err != nil {
			log.Error().Err(errors.New("ERROR QUERY PLTA DELETING : " + err.Error())).Msg("")
			return err
		}

		return nil
	})
}

func (t *pltaRepository) GetByIdWithPreloads(id string, preloads ...string) (*models.Plta, error) {
	var plta *models.Plta

	db := t.db

	for _, p := range preloads {
		db = db.Preload(p)
	}

	if err := db.
		Where("id = ?", id).
		First(&plta).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY GET PLTA BY ID WITH SELECTED FIELDS : " + err.Error())).Msg("")
		return nil, err
	}

	return plta, nil
}

func (u *pltaRepository) GetListWithPaginate(cursor *helpers.Cursor, selectedFields string) ([]*models.Plta, *helpers.CursorPagination, error) {
	db := u.db

	if cursor.Search != "" {
		db = db.Where("LOWER(name) LIKE LOWER(?)", "%"+cursor.Search+"%")
	}

	var total int64
	if err := db.Table("plta").Count(&total).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY PLTA TOTAL : " + err.Error())).Msg("")
		return nil, nil, err
	}

	var sortBy string
	switch strings.ToLower(cursor.SortBy) {
	case "name":
		sortBy = "name"
	case "createdat":
		sortBy = "created_at"
	default:
		sortBy = "created_at"
	}

	var pltas []*models.Plta
	if err := db.Debug().
		Offset((cursor.CurrentPage - 1) * cursor.PerPage).
		Limit(cursor.PerPage).
		Order(fmt.Sprintf("%v %v", sortBy, cursor.SortOrder)).
		Find(&pltas).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY PLTA LIST WITH PAGINATE : " + err.Error())).Msg("")
		return nil, nil, err
	}

	if total == 0 {
		return nil, &helpers.CursorPagination{}, nil
	}

	cursorPagination := cursor.GeneratePager(total)

	return pltas, cursorPagination, nil
}

func (u *pltaRepository) GetTotal() (int64, error) {
	var counter int64

	if err := u.db.
		Table("plta").
		Where("deleted_at IS NULL").
		Count(&counter).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY TOTAL USER : " + err.Error()))
		return 0, err
	}

	return counter, nil
}
