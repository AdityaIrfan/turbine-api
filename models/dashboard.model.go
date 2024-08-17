package models

type DashboardResponse struct {
	TotalUserActive    int64 `json:"TotalUserActive"`
	TotalUserRequest   int64 `json:"TotalUserRequest"`
	TotalPlta          int64 `json:"TotalPlta"`
	TotalTurbineReport int64 `json:"TotalTurbineReport"`
}
