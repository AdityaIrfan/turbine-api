package services

import contract "turbine-api/contracts"

type roleService struct {
	roleRepo contract.IRoleRepository
}

func NewRoleService(roleRepo contract.IRoleRepository) contract.IRoleService {
	return &roleService{
		roleRepo: roleRepo,
	}
}
