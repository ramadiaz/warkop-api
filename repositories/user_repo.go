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
		`INSERT INTO users (email, username, password, first_name, last_name, contact, address) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		data.Email, data.Username, data.Password, data.FirstName, data.LastName, data.Contact, data.Address,
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

	_, err = tx.Exec("INSERT INTO verification_token (user_id, token, expired_at) VALUES($1, $2, NOW() + INTERVAL '2 hours')", data.ID, token)
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

	err := r.DB.QueryRow("SELECT * FROM users WHERE username = $1", username).Scan(&data.ID, &data.Email, &data.Username, &data.Password, &data.FirstName, &data.LastName, &data.Contact, &data.Address, &data.IsVerified, &data.VerifiedAt, &data.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
