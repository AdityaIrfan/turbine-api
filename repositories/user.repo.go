package repositories

import (
	"errors"
	contract "turbine-api/contracts"
	"turbine-api/helpers"
	"turbine-api/models"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) contract.IUserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) IsUsernameExist(username string) (bool, error) {
	var user *models.User

	db := u.db
	if err := db.Select("username").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (u *userRepository) GetById(id string) (*models.User, error) {
	var user *models.User

	db := u.db
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (u *userRepository) GetAllWithPaginate(cursor *helpers.Cursor) ([]*models.User, *helpers.CursorPagination) {
	return nil, nil
}
