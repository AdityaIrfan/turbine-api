package repositories

import (
	"errors"
	contract "pln/AdityaIrfan/turbine-api/contracts"
	"pln/AdityaIrfan/turbine-api/models"
	"strings"

	"github.com/phuslu/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) contract.IRoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) GetAll(search string) ([]*models.Role, error) {
	var roles []*models.Role

	db := r.db
	if search != "" {
		db = db.Where("LOWER(type) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	if err := db.Find(&roles).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY ALL ROLES : " + err.Error())).Msg("")
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) GetById(id string) (*models.Role, error) {
	var role *models.Role

	db := r.db
	if err := db.Where("id = ?", id).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY ROLE BY ID : " + err.Error())).Msg("")
		return nil, err
	}

	return role, nil
}

func (r *roleRepository) GetByIdWithSelectedFields(id string, selectedFields string) (*models.Role, error) {
	var role *models.Role

	db := r.db
	if err := db.Select(selectedFields).Where("id = ?", id).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY ROLE BY ID WITH SELECTED FIELDS : " + err.Error())).Msg("")
		return nil, err
	}

	return role, nil
}

func (r *roleRepository) GetByTypeWithSelectedFields(roleType models.RoleType, selectedFields string) (*models.Role, error) {
	var role *models.Role

	db := r.db
	if err := db.Select(selectedFields).Where("type = ?", roleType).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY ROLE BY TYPE WITH SELECTED FIELDS : " + err.Error())).Msg("")
		return nil, err
	}

	return role, nil
}

func (r *roleRepository) IsEqualTypeExist(roleType models.RoleType) (bool, error) {
	var role *models.Role

	db := r.db

	err := db.Where("LOWER(type) LIKE ?", "%"+strings.ToLower(string(roleType))+"%").First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		log.Error().Err(errors.New("ERROR QUERY ROLE IsEqualTypeExist : " + err.Error())).Msg("")
		return false, err
	}

	return true, nil
}

func (r *roleRepository) Create(role *models.Role) error {
	if err := r.db.Create(&role).Clauses(clause.Returning{}).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY ROLE INSERT : " + err.Error())).Msg("")
		return err
	}

	return nil
}

func (r *roleRepository) Update(role *models.Role) error {
	if err := r.db.Updates(&role).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY ROLE UPDATE : " + err.Error())).Msg("")
		return err
	}

	return nil
}

func (r *roleRepository) Delete(role *models.Role) error {
	if err := r.db.Delete(&role).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY ROLE DELETE : " + err.Error())).Msg("")
		return err
	}

	return nil
}
