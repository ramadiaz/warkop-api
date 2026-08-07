package repositories

import "gorm.io/gorm"

func (r *compRepository) BeginTransaction() (*gorm.DB, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}