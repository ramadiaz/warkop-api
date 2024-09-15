package repositories

import "warkop-api/dto"

func (r *compRepository) RegisterMenu(data dto.Menu) error {
	_, err := r.DB.Exec("INSERT INTO menu (name, type, price, stock) VALUES($1, $2, $3, $4)", data.Name, data.Type, data.Price, data.Stock)
	if err != nil {
		return err
	}

	return nil
}

func (r *compRepository) GetAllMenu() ([]*dto.Menu, error) {
	rows, err := r.DB.Query("SELECT * FROM menu ORDER BY name ASC")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []*dto.Menu

	for rows.Next() {
		var menu dto.Menu
		err = rows.Scan(&menu.ID, &menu.Name, &menu.Type, &menu.Price, &menu.Stock, &menu.CreatedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, &menu)
	}

	return result, nil
}
