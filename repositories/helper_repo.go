package repositories

import "database/sql"

func (r *compRepository) BeginTransaction() (*sql.Tx, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}