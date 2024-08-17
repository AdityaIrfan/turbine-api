package handlers

import (
	contract "pln/AdityaIrfan/turbine-api/contracts"

	"github.com/labstack/echo/v4"
)

type dashboardHandler struct {
	dashboardService contract.IDashboardService
}

func NewDashboardHandler(dashboardService contract.IDashboardService) contract.IDashboardHandler {
	return &dashboardHandler{
		dashboardService: dashboardService,
	}
}

func (d *dashboardHandler) GetDashboardData(c echo.Context) error {
	return d.dashboardService.GetDashboardData(c)
}
