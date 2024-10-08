package contract

import (
	"pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"
	"time"
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
	GetByNameWithSelectedFields(divisionName models.DivisionName, selectedFields string) (*models.Division, error)
	IsEqualNameExist(divisionName models.DivisionName) (bool, error)
	Delete(role *models.Division) error
	GetAllWithPaginate(cursor *helpers.Cursor) ([]*models.Division, *helpers.CursorPagination, error)
}

type IUserRepository interface {
	Create(user *models.User, preloads ...string) error
	Update(user *models.User, preloads ...string) error
	IsUsernameExist(username string) (bool, error)
	IsEmailExist(email string) (bool, error)
	IsPhoneExist(phone string) (bool, error)
	GetById(id string, preloads ...string) (*models.User, error)
	GetByIdWithSelectedFields(id string, selectedFields string, preloads ...string) (*models.User, error)
	GetByUsernameWithSelectedFields(username string, selectedFields string, preloads ...string) (*models.User, error)
	GetByEmailWithSelectedFields(email string, selectedFields string, preloads ...string) (*models.User, error)
	GetByPhoneWithSelectedFields(phone string, selectedFields string, preloads ...string) (*models.User, error)
	GetAllWithPaginate(cursor *helpers.Cursor, userRoles ...models.UserRole) ([]*models.User, *helpers.CursorPagination, error)
	Delete(user *models.User) error
	GetTotalByStatus(status models.UserStatus) (int64, error)
}

type IAuthRedisRepository interface {
	SaveRefreshToken(id string, refreshToken *models.RefreshTokenRedis, ttl time.Duration)
	GetRefreshToken(id string) (*models.RefreshTokenRedis, error)
	DeleteRefreshToken(id string)
	IncLoginFailedCounter(id string)
	IsLoginBlocked(id string) (bool, error)
	SaveToken(id string, token string, ttl time.Duration) error
	GetToken(id string) (string, error)
	DeleteToken(id string)
}

type IConfigRepository interface {
	SaveOrUpdateRootLocation(rootLocation *models.ConfigRootLocation) error
	GetByType(configType models.ConfigType) (*models.Config, error)
}

type IConfigRedisRepository interface {
	SaveRootLocation(rootLocation *models.ConfigRootLocation)
	GetRootLocation() (*models.ConfigRootLocation, error)
}

type ITurbineRepository interface {
	Create(turbine *models.Turbine) error
	GetByIdWithSelectedFields(id string, selectedFields string, preloads ...string) (*models.Turbine, error)
	GetAllWithPaginate(cursor *helpers.Cursor, selectedFields string) ([]*models.Turbine, *helpers.CursorPagination, error)
	GetLatest() (*models.Turbine, error)
	Delete(turbine *models.Turbine) error
	GetTotal() (int64, error)
}

type IPltaRepository interface {
	Create(plta *models.Plta) error
	Update(plta *models.Plta) error
	GetByIdWithSelectedFields(id string, selectedFields string) (*models.Plta, error)
	GetByEqualNameWithSelectedFields(name string, selectedFields string) (*models.Plta, error)
	GetByIdWithPreloads(id string, preloads ...string) (*models.Plta, error)
	GetAll(search string) ([]*models.Plta, error)
	Delete(plta *models.Plta) error
	GetListWithPaginate(cursor *helpers.Cursor, selectedFields string) ([]*models.Plta, *helpers.CursorPagination, error)
	GetTotal() (int64, error)
}

type IPltaUnitRepository interface {
	GetByIdAndSelectedFields(id, selectedFields string) (*models.PltaUnit, error)
	CreateOrUpdate(pltaUnits []*models.PltaUnit) ([]*models.PltaUnit, error)
	Delete(pltaUnit *models.PltaUnit) error
	GetByIdWithPreloads(id string, preloads ...string) (*models.PltaUnit, error)
	GetAll(search string) ([]*models.PltaUnit, error)
}
