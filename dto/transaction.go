package dto

type Transaction struct {
	ID        int64              `json:"id"`
	CashierID string             `json:"cashier_id"`
	Cashier   string             `json:"cashier"`
	Total     int64                `json:"total"`
	Cash      int64              `json:"cash"`
	Change    int64              `json:"change"`
	Menus     []*TransactionItem `json:"menus"`
	CreatedAt string             `json:"created_at"`
}

type TransactionItem struct {
	ID            int    `json:"id"`
	TransactionID int64  `json:"transaction_id"`
	MenuID        int    `json:"menu_id"`
	Name          string `json:"name"`
	Price         int64  `json:"price"`
	Amount        int64  `json:"amount"`
	Quantity      int    `json:"quantity"`
	CreatedAt     string `json:"created_at"`
}
