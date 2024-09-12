package services

import (
	"warkop-api/dto"
	"warkop-api/repositories"
)

type CompService interface {
	RegisterUser(data dto.User) (*string, error)
	GenerateJWT(data dto.User) (*string, error)
}

type compServices struct {
	repo repositories.CompRepository
}

func NewService(r repositories.CompRepository) *compServices {
	return &compServices{
		repo: r,
	}
}
