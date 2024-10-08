package repositories

import (
	"database/sql"
	"warkop-api/dto"
)

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
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`INSERT INTO transaction_item (transaction_id, menu_id, quantity) VALUES($1, $2, $3)`,
		data.TransactionID, data.MenuID, data.Quantity,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		UPDATE menu SET stock = (SELECT stock FROM menu WHERE id = $1) - $2 WHERE id = $1
	`, data.MenuID, data.Quantity)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *compRepository) GetTransaction(id string) (*dto.Transaction, error) {
	var data dto.Transaction

	err := r.DB.QueryRow(`
		SELECT transaction.*, users.username 
		FROM transaction 
		JOIN users ON users.id = transaction.cashier_id::uuid 
		WHERE transaction.id = $1;
	`, id).Scan(&data.ID, &data.CashierID, &data.Total, &data.Cash, &data.CreatedAt, &data.Cashier)
	if err != nil {
		return nil, err
	}

	data.Change = data.Cash - data.Total

	return &data, nil
}

func (r *compRepository) GetTransactionItem(id string) ([]*dto.TransactionItem, error) {
	var data []*dto.TransactionItem

	rows, err := r.DB.Query(`
		SELECT transaction_item.*, menu.name, menu.price 
		FROM transaction_item 
		JOIN menu ON menu.id = transaction_item.menu_id 
		WHERE transaction_item.transaction_id = $1;
	`, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var item dto.TransactionItem

		err := rows.Scan(&item.ID, &item.TransactionID, &item.MenuID, &item.Quantity, &item.CreatedAt, &item.Name, &item.Price)
		if err != nil {
			return nil, err
		}

		item.Amount = int64(item.Quantity) * item.Price

		data = append(data, &item)
	}

	return data, nil
}

func (r *compRepository) GetTransactionItemInTx(tx *sql.Tx, id string) ([]*dto.TransactionItem, error) {
	var data []*dto.TransactionItem

	rows, err := tx.Query(`
		SELECT transaction_item.*, menu.name, menu.price 
		FROM transaction_item 
		JOIN menu ON menu.id = transaction_item.menu_id 
		WHERE transaction_item.transaction_id = $1;
	`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item dto.TransactionItem
		err := rows.Scan(&item.ID, &item.TransactionID, &item.MenuID, &item.Quantity, &item.CreatedAt, &item.Name, &item.Price)
		if err != nil {
			return nil, err
		}
		item.Amount = int64(item.Quantity) * item.Price
		data = append(data, &item)
	}

	return data, nil
}

func (r *compRepository) GetAllTransaction() ([]*dto.Transaction, error) {
	var data []*dto.Transaction

	rows, err := r.DB.Query(`
		SELECT transaction.*, users.username 
		FROM transaction 
		JOIN users ON users.id = transaction.cashier_id::uuid 
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var tx dto.Transaction

		err := rows.Scan(&tx.ID, &tx.CashierID, &tx.Total, &tx.Cash, &tx.CreatedAt, &tx.Cashier)
		if err != nil {
			return nil, err
		}

		tx.Change = tx.Cash - tx.Total

		data = append(data, &tx)
	}

	return data, nil
}
