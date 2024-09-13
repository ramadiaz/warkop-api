package services

import (
	"errors"
	"os"
	"warkop-api/helpers"
)

func (s *compServices) GenerateAPIKey(name string, secret string) (*string, error) {
	admin_secret := os.Getenv("ADMIN_SECRET")
	if secret != admin_secret {
		return nil, errors.New("invalid secret")
	}

	token, err := helpers.GenerateToken(64)
	if err != nil {
		return nil, err
	}

	err = s.repo.RegisterAPIKey(name, token)
	if err != nil {
		return nil, err
	}

	return &token, nil

}
