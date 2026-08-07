package repositories

import (
	"warkop-api/dto"
	"warkop-api/models"
)

func (r *compRepository) RegisterMenu(data dto.Menu) error {
	menu := models.Menu{
		Name:  data.Name,
		Type:  data.Type,
		Price: data.Price,
		Stock: data.Stock,
	}
	return r.DB.Create(&menu).Error
}

func (r *compRepository) GetAllMenu() ([]*dto.Menu, error) {
	var dbMenus []models.Menu
	err := r.DB.Order("name ASC").Find(&dbMenus).Error
	if err != nil {
		return nil, err
	}

	var result []*dto.Menu
	for _, menu := range dbMenus {
		result = append(result, &dto.Menu{
			ID:        int(menu.ID),
			Name:      menu.Name,
			Type:      menu.Type,
			Price:     menu.Price,
			Stock:     menu.Stock,
			CreatedAt: menu.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}
