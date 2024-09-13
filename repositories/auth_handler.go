package repositories

func (r *compRepository) RegisterAPIKey(name string, secret string) error {
	_, err := r.DB.Exec("INSERT INTO api_key (name, token) VALUES($1, $2)", name, secret)
	if err != nil {
		return err
	}

	return nil
}
