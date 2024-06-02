package contract

import (
	"time"
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
	Create(user *models.User, preloads ...string) error
	Update(user *models.User, preloads ...string) error
	IsUsernameExist(username string) (bool, error)
	GetById(id string, preloads ...string) (*models.User, error)
	GetByIdWithSelectedFields(id string, selectedFields string) (*models.User, error)
	GetByUsernameWithSelectedFields(username string, selectedFields string, preloads ...string) (*models.User, error)
	GetAllWithPaginate(cursor *helpers.Cursor) ([]*models.User, *helpers.CursorPagination, error)
	Delete(user *models.User) error
}

type IAuthRedisRepository interface {
	SaveRefreshToken(id string, refreshToken *models.RefreshTokenRedis, ttl time.Duration)
	GetRefreshToken(id string) (*models.RefreshTokenRedis, error)
	DeleteRefreshToken(id string)
	IncLoginFailedCounter(id string)
	IsLoginBlocked(id string) (bool, error)
}

type IConfigRepository interface {
	GetByType(configType models.ConfigType) (*models.Config, error)
}

type IConfigRedisRepository interface {
	SaveRootLocation(rootLocation *models.ConfigRootLocation)
	GetRootLocation() (*models.ConfigRootLocation, error)
}

type ITurbineRepository interface{}
