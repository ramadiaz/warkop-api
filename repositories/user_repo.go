package repositories

import "warkop-api/dto"

func (r *compRepository) RegisterUser(data dto.User) (int64, error) {
	var id int64
	err := r.DB.QueryRow(
		`INSERT INTO users (email, username, password, first_name, last_name, contact, address) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		data.Email, data.Username, data.Password, data.FirstName, data.LastName, data.Contact, data.Address,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
