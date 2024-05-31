package repositories

import (
	"errors"
	"strings"
	contract "turbine-api/contracts"
	"turbine-api/models"

	"gorm.io/gorm"
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
		db = db.Where("LOWER(type) LIKES ?", "'%"+strings.ToLower(search)+"%'")
	}

	err := db.Find(&roles).Error
	return roles, err
}

func (r *roleRepository) GetById(id string) (*models.Role, error) {
	var role *models.Role

	db := r.db
	if err := db.Where("id = ?", id).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return role, nil
}

func (r *roleRepository) IsEqualTypeExist(roleType models.RoleType) (bool, error) {
	var role *models.Role

	db := r.db

	err := db.Where("LOWER(type) LIKES ?", "'%"+strings.ToLower(string(roleType))+"%'").First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
