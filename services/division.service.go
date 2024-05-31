package services

import contract "turbine-api/contracts"

type divisionService struct {
	divisionRepo contract.IDivisionRepository
}

func NewDivisionService(divisionRepo contract.IDivisionRepository) contract.IDivisionService {
	return &divisionService{
		divisionRepo: divisionRepo,
	}
}
