package repositories

import "warkop-api/dto"

func (r *compRepository) RegisterMenu(data dto.Menu) error {
	_, err := r.DB.Exec("INSERT INTO menu (name, type, price, stock) VALUES($1, $2, $3, $4)", data.Name, data.Type, data.Price, data.Stock)
	if err != nil {
		return err
	}

	return nil
}
