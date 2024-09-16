package services

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"warkop-api/dto"
	"warkop-api/helpers"

	"github.com/dgrijalva/jwt-go"
)

func (s *compServices) RegisterUser(data dto.User) error {
	id, err := s.repo.RegisterUser(data)
	if err != nil {
		return err
	}

	data.ID = *id
	data.IsVerified = false

	err = s.GenerateVerificationEmail(data.Username)
	if err != nil {
		return err
	}

	return nil
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

func (s *compServices) VerifyAccount(token string) error {
	return s.repo.VerifyAccount(token)
}

func (s *compServices) LoginUser(username string, password string) (*string, error) {
	data, err := s.repo.GetUserData(username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(strconv.Itoa(http.StatusNotFound))
		}
		return nil, err
	}

	if password != data.Password {
		return nil, errors.New(strconv.Itoa(http.StatusUnauthorized))
	}

	if !data.IsVerified {
		return nil, errors.New(strconv.Itoa(http.StatusForbidden))
	}

	token, err := s.GenerateJWT(*data)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *compServices) RequestResetPassword(username string) (*dto.User, error) {
	user_data, err := s.repo.GetUserData(username)
	if err != nil {
		return nil, err
	}

	otp, err := helpers.GenerateNumericOTP()
	if err != nil {
		return nil, err
	}

	err = s.repo.RequestResetPassword(*user_data, otp)
	if err != nil {
		return nil, err
	}

	err = s.GenerateOTPEmail(*user_data, otp)
	if err != nil {
		return nil, err
	}

	return user_data, nil
}

func (s *compServices) VerifyResetPassword(data dto.OTPVerifyToken) (*string, error) {
	otp_data, err := s.repo.VerifyResetPassword(data)
	if err != nil {
		return nil, err
	}

	if otp_data.OTP != data.OTP {
		return nil, errors.New(strconv.Itoa(http.StatusUnauthorized))
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = data.UserID
	claims["otp"] = data.OTP

	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	secret := os.Getenv("JWT_SECRET")

	secretKey := []byte(secret)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (s *compServices) ResetPassword(user_data dto.User) error {
	return s.repo.ResetPassword(user_data)
}

func (s *compServices) UploadUserProfile(data dto.User, image_url string) error {
	return s.repo.UploadUserProfile(data, image_url)
}

func (s *compServices) GetUserProfile(id string) (*string, error) {
	return s.repo.GetUserProfile(id)
}

func (s *compServices) GenerateVerificationEmail(username string) error {
	base_url := os.Getenv("FRONT_END_BASE_URL")

	data, err := s.repo.GetUserData(username)
	if err != nil {
		return err
	}

	if data.IsVerified {
		return errors.New("account already verified")
	}

	token, err := s.repo.RegisterToken(*data)
	if err != nil {
		return err
	}

	verify_url := fmt.Sprintf(base_url+"/verify?token=%s", *token)

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
					<a href="%s" class="message" style="font-size: 24px; font-weight: bold;">Verify Now</a>
					<p class="message">This code will expire in 2 hours. If you did not request this verification, please ignore this email.</p>
					<p class="footer">Best regards,<br>Wyvern Team</p>
				</div>
			</body>
		</html>`,
		data.Email, verify_url,
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

func (s *compServices) GenerateOTPEmail(data dto.User, otp string) error {
	body := fmt.Sprintf(
		`<html>
			<head>
				<title>OTP Code Verification</title>
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
					<p class="title">OTP Code Verification</p>
					<p class="message">Dear %s,</p>
					<p class="message">We received a request to reset your password for your Warkop Cashier account. If you didn’t make this request, please ignore this email. Otherwise, you can reset your password by input the OTP Code below.</p>
					<h2 class="message" style="font-size: 24px; font-weight: bold;">%s</h2>
					<p class="message">This code will expire in 2 hours. If you did not request this verification, please ignore this email.</p>

					<p class="footer">Best regards,<br>Wyvern Team</p>
				</div>
			</body>
		</html>`,
		data.Email, otp,
	)

	var email_data dto.Email

	email_data.Email = data.Email
	email_data.Subject = "Warkop Cashier System - OTP Code"
	email_data.Body = body

	err := s.SendEmail(email_data)
	if err != nil {
		return err
	}

	return nil
}
