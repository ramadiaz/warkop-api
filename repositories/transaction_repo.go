package repositories

import (
	"errors"
	"warkop-api/dto"
	"warkop-api/models"

	"gorm.io/gorm"
)

func (r *compRepository) RegisterTransaction(data dto.Transaction) (*int64, error) {
	tx := models.Transaction{
		CashierID: data.CashierID,
		Total:     data.Total,
		Cash:      data.Cash,
	}

	err := r.DB.Create(&tx).Error
	if err != nil {
		return nil, err
	}

	return &tx.ID, nil
}

func (r *compRepository) RegisterTransactionItem(data dto.TransactionItem) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		item := models.TransactionItem{
			TransactionID: data.TransactionID,
			MenuID:        data.MenuID,
			Quantity:      data.Quantity,
		}

		err := tx.Create(&item).Error
		if err != nil {
			return err
		}

		var menu models.Menu
		err = tx.Where("id = ?", data.MenuID).First(&menu).Error
		if err != nil {
			return err
		}

		newStock := menu.Stock - data.Quantity
		return tx.Model(&menu).Update("stock", newStock).Error
	})
}

func (r *compRepository) GetTransaction(id string) (*dto.Transaction, error) {
	var tx models.Transaction

	err := r.DB.Where("id = ?", id).First(&tx).Error
	if err != nil {
		return nil, err
	}

	var cashierUsername string
	err = r.DB.Model(&models.User{}).Select("username").Where("id = ?", tx.CashierID).Scan(&cashierUsername).Error
	if err != nil {
		return nil, err
	}

	data := dto.Transaction{
		ID:        tx.ID,
		CashierID: tx.CashierID,
		Cashier:   cashierUsername,
		Total:     tx.Total,
		Cash:      tx.Cash,
		Change:    tx.Cash - tx.Total,
		CreatedAt: tx.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return &data, nil
}

func (r *compRepository) GetTransactionItem(id string) ([]*dto.TransactionItem, error) {
	var items []models.TransactionItem

	err := r.DB.Where("transaction_id = ?", id).Find(&items).Error
	if err != nil {
		return nil, err
	}

	var data []*dto.TransactionItem
	for _, item := range items {
		var menu models.Menu
		err = r.DB.Where("id = ?", item.MenuID).First(&menu).Error
		if err != nil {
			return nil, err
		}

		data = append(data, &dto.TransactionItem{
			ID:            int(item.ID),
			TransactionID: item.TransactionID,
			MenuID:        item.MenuID,
			Name:          menu.Name,
			Price:         int64(menu.Price),
			Amount:        int64(item.Quantity) * int64(menu.Price),
			Quantity:      item.Quantity,
			CreatedAt:     item.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return data, nil
}

func (r *compRepository) GetTransactionItemInTx(tx *gorm.DB, id string) ([]*dto.TransactionItem, error) {
	var items []models.TransactionItem

	err := tx.Where("transaction_id = ?", id).Find(&items).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var data []*dto.TransactionItem
	for _, item := range items {
		var menu models.Menu
		err = tx.Where("id = ?", item.MenuID).First(&menu).Error
		if err != nil {
			return nil, err
		}

		data = append(data, &dto.TransactionItem{
			ID:            int(item.ID),
			TransactionID: item.TransactionID,
			MenuID:        item.MenuID,
			Name:          menu.Name,
			Price:         int64(menu.Price),
			Amount:        int64(item.Quantity) * int64(menu.Price),
			Quantity:      item.Quantity,
			CreatedAt:     item.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return data, nil
}

func (r *compRepository) GetAllTransaction() ([]*dto.Transaction, error) {
	var transactions []models.Transaction

	err := r.DB.Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	var data []*dto.Transaction
	for _, tx := range transactions {
		var cashierUsername string
		err = r.DB.Model(&models.User{}).Select("username").Where("id = ?", tx.CashierID).Scan(&cashierUsername).Error
		if err != nil {
			return nil, err
		}

		data = append(data, &dto.Transaction{
			ID:        tx.ID,
			CashierID: tx.CashierID,
			Cashier:   cashierUsername,
			Total:     tx.Total,
			Cash:      tx.Cash,
			Change:    tx.Cash - tx.Total,
			CreatedAt: tx.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return data, nil
}
