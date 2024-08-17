package services

import (
	"net/http"
	contract "pln/AdityaIrfan/turbine-api/contracts"
	"pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"

	"github.com/labstack/echo/v4"
)

type dashbaordServices struct {
	userRepo    contract.IUserRepository
	turbineRepo contract.ITurbineRepository
	pltaRepo    contract.IPltaRepository
}

func NewDashboardService(
	userRepo contract.IUserRepository,
	turbineRepo contract.ITurbineRepository,
	pltaRepo contract.IPltaRepository,
) contract.IDashboardService {
	return &dashbaordServices{
		userRepo:    userRepo,
		turbineRepo: turbineRepo,
		pltaRepo:    pltaRepo,
	}
}

func (d *dashbaordServices) GetDashboardData(c echo.Context) error {
	totalUserActive, err := d.userRepo.GetTotalByStatus(models.UserStatus_Active)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	totalUserRequest, err := d.userRepo.GetTotalByStatus(models.UserStatus_InActive)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	totalPlta, err := d.pltaRepo.GetTotal()
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	totalTurbineReport, err := d.turbineRepo.GetTotal()
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan data dashboard", &models.DashboardResponse{
		TotalUserActive:    totalUserActive,
		TotalUserRequest:   totalUserRequest,
		TotalPlta:          totalPlta,
		TotalTurbineReport: totalTurbineReport,
	})
}
