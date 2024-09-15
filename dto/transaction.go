package dto

type Transaction struct {
	ID        int64             `json:"id"`
	CashierID string            `json:"cashier_id"`
	Total     int               `json:"total"`
	Cash      int               `json:"cash"`
	Menus     []TransactionItem `json:"menus"`
	CreatedAt string            `json:"created_at"`
}

type TransactionItem struct {
	ID            int    `json:"id"`
	TransactionID int64  `json:"transaction_id"`
	MenuID        int    `json:"menu_id"`
	Quantity      int    `json:"quantity"`
	CreatedAt     string `json:"created_at"`
}
