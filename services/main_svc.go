package services

import (
	"warkop-api/dto"
	"warkop-api/repositories"
)

type CompService interface {
	RegisterUser(data dto.User) error
	GenerateJWT(data dto.User) (*string, error)
	SendEmail(data dto.Email) error
	GenerateVerificationEmail(username string) error
	GenerateAPIKey(name string, secret string) (*string, error)
	VerifyAccount(token string) error
	LoginUser(username string, password string) (*string, error)
}

type compServices struct {
	repo repositories.CompRepository
}

func NewService(r repositories.CompRepository) *compServices {
	return &compServices{
		repo: r,
	}
}
