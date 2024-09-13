package services

import (
	"fmt"
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

	err = s.GenerateVerificationEmail(data.Username)
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

func (s *compServices) GenerateVerificationEmail(username string) error {
	base_url := os.Getenv("FRONT_END_BASE_URL")

	data, err := s.repo.GetUserData(username)
	if err != nil {
		return err
	}

	token, err := s.repo.RegisterToken(*data)
	if err != nil {
		return err
	}

	body := fmt.Sprintf(
		`<html>
			<head>
				<title>Email Verification</title>
				<style>
					body {
						font-family: Arial, sans-serif;
						margin: 0;
						padding: 0;
					}
					.container {
						max-width: 600px;
						margin: 20px auto;
						padding: 20px;
						border: 1px solid #ccc;
						border-radius: 5px;
						background-color: #f9f9f9;
					}
					.title {
						font-size: 24px;
						font-weight: bold;
						margin-bottom: 20px;
					}
					.message {
						margin-bottom: 20px;
					}
					.footer {
						margin-top: 20px;
						font-size: 14px;
						color: #666;
					}
				</style>
			</head>
			<body>
				<div class="container">
					<p class="title">Email Verification</p>
					<p class="message">Dear %s,</p>
					<p class="message">Thank you for registering with our platform. To complete your registration, please clik the button below</p>
					<a href="%s?token=%s" class="message" style="font-size: 24px; font-weight: bold;">Verify Now</a>
					<p class="message">This code will expire in 24 hours. If you did not request this verification, please ignore this email.</p>
					<p class="footer">Best regards,<br>Wyvern Team</p>
				</div>
			</body>
		</html>`,
		data.Email, base_url, *token,
	)

	var email_data dto.Email

	email_data.Email = data.Email
	email_data.Subject = "Warkop Cashier System - Verification"
	email_data.Body = body

	err = s.SendEmail(email_data)
	if err != nil {
		return err
	}

	return nil
}

func (s *compServices) VerifyAccount(token string) error {
	return s.repo.VerifyAccount(token)
}
