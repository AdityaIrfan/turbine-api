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

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) contract.IUserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Create(user *models.User, preloads ...string) error {
	db := u.db

	for _, p := range preloads {
		db = db.Preload(p)
	}

	if err := db.Create(&user).Clauses(clause.Returning{}).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY USER CREATE : " + err.Error())).Msg("")
		return err
	}

	return nil
}

func (u *userRepository) Update(user *models.User, preloads ...string) error {

	db := u.db

	for _, p := range preloads {
		db = db.Preload(p)
	}
	if err := db.Updates(&user).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY USER UPDATE : " + err.Error())).Msg("")
		return err
	}

	return nil
}

func (u *userRepository) IsUsernameExist(username string) (bool, error) {
	var user *models.User

	if err := u.db.Select("username").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		log.Error().Err(errors.New("ERROR QUERY IsUsernameExist : " + err.Error())).Msg("")
		return false, err
	}

	return true, nil
}

func (u *userRepository) IsEmailExist(email string) (bool, error) {
	var user *models.User

	if err := u.db.Select("email").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		log.Error().Err(errors.New("ERROR QUERY IsEmailExist : " + err.Error())).Msg("")
		return false, err
	}

	return true, nil
}

func (u *userRepository) GetById(id string, preloads ...string) (*models.User, error) {
	var user *models.User

	db := u.db

	for _, p := range preloads {
		db = db.Preload(p)
	}

	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY USER BY ID : " + err.Error())).Msg("")
		return nil, err
	}

	return user, nil
}

func (u *userRepository) GetByIdWithSelectedFields(id string, selectedFields string, preloads ...string) (*models.User, error) {
	var user *models.User

	db := u.db

	for _, p := range preloads {
		db = db.Preload(p)
	}

	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY USER BY ID WITH SELECTED FIELDS: " + err.Error())).Msg("")
		return nil, err
	}

	return user, nil
}

func (u *userRepository) GetByUsernameWithSelectedFields(username string, selectedFields string, preloads ...string) (*models.User, error) {
	var user *models.User

	db := u.db

	for _, p := range preloads {
		db = db.Preload(p)
	}

	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY USER BY USERNAME WITH SELECTED FIELDS : " + err.Error())).Msg("")
		return nil, err
	}

	return user, nil
}

func (u *userRepository) GetByEmailWithSelectedFields(email string, selectedFields string, preloads ...string) (*models.User, error) {
	var user *models.User

	db := u.db

	for _, p := range preloads {
		db = db.Preload(p)
	}

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY USER BY EMAIL WITH SELECTED FIELDS : " + err.Error())).Msg("")
		return nil, err
	}

	return user, nil
}

func (u *userRepository) GetAllWithPaginate(cursor *helpers.Cursor) ([]*models.User, *helpers.CursorPagination, error) {
	db := u.db

	if cursor.Search != "" {
		db = db.Where("LOWER(name) LIKE LOWER(?)", "%"+cursor.Search+"%")
	}

	var total int64
	if err := db.Table("users").Count(&total).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY USERS TOTAL : " + err.Error())).Msg("")
		return nil, nil, err
	}

	var sortBy string
	switch strings.ToLower(cursor.SortBy) {
	case "name":
		sortBy = "name"
	case "createdat":
		sortBy = "created_at"
	case "username":
		sortBy = "username"
	default:
		sortBy = "created_at"
	}

	var users []*models.User
	if err := db.
		Offset(cursor.CurrentPage - 1).
		Limit(cursor.PerPage).
		Order(fmt.Sprintf("%v %v", sortBy, cursor.SortOrder)).
		Find(&users).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY USERS LIST WITH PAGINATE : " + err.Error())).Msg("")
		return nil, nil, err
	}

	cursorPagination := cursor.GeneratePager(total)

	return users, cursorPagination, nil
}

func (u *userRepository) Delete(user *models.User) error {
	if err := u.db.Delete(&user).Error; err != nil {
		log.Error().Err(errors.New("ERROR QUERY USER DELETE : " + err.Error())).Msg("")
		return err
	}

	return nil
}
