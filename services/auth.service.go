package services

import (
	contract "turbine-api/contracts"
	"turbine-api/models"

	"github.com/labstack/echo/v4"
)

type authService struct {
	userRepo contract.IUserRepository
}

func NewAuthService(userRepo contract.IUserRepository) contract.IAuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (a *authService) Register(c echo.Context, in *models.Register) error {
	return nil
}

func (a *authService) Login(c echo.Context, in *models.Login) error {
	return nil
}

func (a *authService) RefreshToken(c echo.Context, in *models.RefreshToken) error {
	return nil
}
