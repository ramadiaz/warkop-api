package repositories

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"warkop-api/dto"
	"warkop-api/helpers"
)

func (r *compRepository) RegisterUser(data dto.User) (*string, error) {
	var id string
	err := r.DB.QueryRow(
		`INSERT INTO users (username, email, password, first_name, last_name, contact, address) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		data.Username, data.Email, data.Password, data.FirstName, data.LastName, data.Contact, data.Address,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *compRepository) RegisterToken(data dto.User) (*string, error) {
	token, err := helpers.GenerateToken(32)
	if err != nil {
		return nil, err
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("DELETE FROM verification_token WHERE user_id = $1", data.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO verification_token (user_id, token, expired_at) VALUES($1, $2, NOW() + INTERVAL '9 hours')", data.ID, token)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &token, nil
}

func (r *compRepository) VerifyAccount(token string) error {
	var expiredAt time.Time

	err := r.DB.QueryRow("SELECT expired_at FROM verification_token WHERE token = $1", token).Scan(&expiredAt)
	if err != nil {
		return err
	}

	if time.Now().After(expiredAt) {
		return errors.New(strconv.Itoa(http.StatusGone))
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE users SET is_verified = true, verified_at = NOW() WHERE id = (SELECT user_id FROM verification_token WHERE token = $1)", token)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`DELETE from verification_token WHERE token = $1`, token)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *compRepository) GetUserData(username string) (*dto.User, error) {
	var data dto.User

	err := r.DB.QueryRow("SELECT * FROM users WHERE username = $1", username).Scan(&data.ID, &data.Username, &data.Email, &data.Password, &data.FirstName, &data.LastName, &data.Contact, &data.Address, &data.IsVerified, &data.VerifiedAt, &data.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *compRepository) RequestResetPassword(data dto.User, otp string) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM reset_otp WHERE user_id = $1::uuid
	`, data.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO reset_otp (user_id, otp, expired_at) VALUES($1, $2, NOW() + INTERVAL '2 hours')
	`, data.ID, otp)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *compRepository) VerifyResetPassword(data dto.OTPVerifyToken) (*dto.OTPVerifyToken, error) {
	var d dto.OTPVerifyToken

	err := r.DB.QueryRow(`
		SELECT user_id, otp FROM reset_otp WHERE user_id = $1
	`, data.UserID).Scan(&d.UserID, &d.OTP)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (r *compRepository) ResetPassword(user_data dto.User) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM reset_otp WHERE user_id = $1::uuid", user_data.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE users SET password = $1 WHERE id = $2::uuid", user_data.Password, user_data.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *compRepository) UploadUserProfile(data dto.User, image_url string) error {
	_, err := r.DB.Exec(`
		INSERT INTO users_image (user_id, image_url)
		VALUES ($1, $2)
		ON CONFLICT (user_id)
		DO UPDATE SET
			image_url = $2;
	`, data.ID, image_url)
	if err != nil {
		return err
	}

	return nil
}

func (r *compRepository) GetUserProfile(id string) (*string, error) {
	var image_url string

	err := r.DB.QueryRow(`
		SELECT image_url FROM users_image WHERE user_id = $1
	`, id).Scan(&image_url)
	if err != nil {
		return nil, err
	}

	return &image_url, nil
}
