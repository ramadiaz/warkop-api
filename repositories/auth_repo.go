package repositories

import "warkop-api/models"

func (r *compRepository) RegisterAPIKey(name string, secret string) error {
	apiKey := models.APIKey{
		Name:  name,
		Token: secret,
	}
	return r.DB.Create(&apiKey).Error
}
