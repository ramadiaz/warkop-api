package repositories

import "warkop-api/dto"

func (r *compRepository) RegisterTransaction(data dto.Transaction) (*int64, error) {
	var id int64
	err := r.DB.QueryRow(
		`INSERT INTO transaction (cashier_id, total, cash) VALUES($1, $2, $3) RETURNING id`, data.CashierID, data.Total, data.Cash,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *compRepository) RegisterTransactionItem(data dto.TransactionItem) error {
	_, err := r.DB.Exec(
		`INSERT INTO transaction_item (transaction_id, menu_id, quantity) VALUES($1, $2, $3)`,
		data.TransactionID, data.MenuID, data.Quantity,
	)
	if err != nil {
		return err
	}

	return nil
}
