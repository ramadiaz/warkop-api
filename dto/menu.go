package dto

type Menu struct {
	ID        int    `json:"id"`
	Name      string `json:"name" form:"name" binding:"required"`
	Type      string `json:"type" form:"type" binding:"required"`
	Price     int    `json:"price" form:"price" binding:"required"`
	Stock     int    `json:"stock" form:"stock" binding:"required"`
	CreatedAt string `json:"created_at"`
}