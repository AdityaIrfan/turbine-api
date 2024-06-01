package contract

import (
	"turbine-api/helpers"
	"turbine-api/models"
)

type IRoleRepository interface {
	Create(role *models.Role) error
	Update(role *models.Role) error
	GetAll(search string) ([]*models.Role, error)
	GetById(id string) (*models.Role, error)
	GetByIdWithSelectedFields(id string, selectedFields string) (*models.Role, error)
	GetByTypeWithSelectedFields(roleType models.RoleType, selectedFields string) (*models.Role, error)
	IsEqualTypeExist(roleType models.RoleType) (bool, error)
	Delete(role *models.Role) error
}

type IDivisionRepository interface {
	Create(division *models.Division) error
	Update(division *models.Division) error
	GetAll(search string) ([]*models.Division, error)
	GetById(id string) (*models.Division, error)
	GetByIdWithSelectedFields(id string, selectedFields string) (*models.Division, error)
	GetByTypeWithSelectedFields(divisionTYpe models.DivisionType, selectedFields string) (*models.Division, error)
	IsEqualTypeExist(divisionType models.DivisionType) (bool, error)
	Delete(role *models.Division) error
}

type IUserRepository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	IsUsernameExist(username string) (bool, error)
	GetById(id string) (*models.User, error)
	GetAllWithPaginate(cursor *helpers.Cursor) ([]*models.User, *helpers.CursorPagination, error)
}

type ITurbineRepository interface{}
