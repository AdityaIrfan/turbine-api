package services

import contract "turbine-api/contracts"

type userService struct {
	userRepo contract.IUserRepository
}

func NewUserService(userRepo contract.IUserRepository) contract.IUserService {
	return &userService{
		userRepo: userRepo,
	}
}
