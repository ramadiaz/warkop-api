package services

import (
	"os"
	"time"
	"warkop-api/dto"

	"github.com/dgrijalva/jwt-go"
)

func (s *compServices) RegisterUser(data dto.User) (*string, error) {
	id, err := s.repo.RegisterUser(data)
	if err != nil {
		return nil, err
	}

	data.ID = id
	data.IsVerified = false

	token, err := s.GenerateJWT(data)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *compServices) GenerateJWT(data dto.User) (*string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = data.ID
	claims["email"] = data.Email
	claims["username"] = data.Username
	claims["first_name"] = data.FirstName
	claims["last_name"] = data.LastName
	claims["contact"] = data.Contact
	claims["address"] = data.Address
	claims["is_verified"] = data.IsVerified

	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	secret := os.Getenv("JWT_SECRET")

	secretKey := []byte(secret)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}
