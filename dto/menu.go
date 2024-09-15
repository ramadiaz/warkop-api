package dto

type Menu struct {
	ID        int    `json:"id"`
	Name      string `json:"name" form:"name" binding:"required"`
	Type      string `json:"type" form:"type" binding:"required"`
	Price     int    `json:"price" form:"price" binding:"required"`
	Stock     int    `json:"stock" form:"stock" binding:"required"`
	CreatedAt string `json:"created_at"`
}

type Transaction struct {
	ID        int    `json:"id"`
	CashierID string `json:"cashier_id"`
	Total     int    `json:"total"`
	Cash      int    `json:"cash"`
	CreatedAt string `json:"created_at"`
}

type TransactionItem struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	MenuID        int    `json:"menu_id"`
	Quantity      int    `json:"quantity"`
	CreatedAt     string `json:"created_at"`
}
