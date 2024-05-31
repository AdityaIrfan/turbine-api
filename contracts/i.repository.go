package contract

import (
	"turbine-api/helpers"
	"turbine-api/models"
)

type IRoleRepository interface {
	GetAll(search string) ([]*models.Role, error)
	GetById(id string) (*models.Role, error)
	IsEqualTypeExist(roleType models.RoleType) (bool, error)
}

type IDivisionRepository interface {
	GetAll(search string) ([]*models.Division, error)
	GetById(id string) (*models.Division, error)
	IsEqualTypeExist(divisionType models.DivisionType) (bool, error)
}

type IUserRepository interface {
	IsUsernameExist(username string) (bool, error)
	GetById(id string) (*models.User, error)
	GetAllWithPaginate(cursor *helpers.Cursor) ([]*models.User, *helpers.CursorPagination)
}
