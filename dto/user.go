package dto

type User struct {
	ID         string  `json:"id"`
	Username   string  `json:"username" form:"username" binding:"required"`
	Email      string  `json:"email" form:"email" binding:"required"`
	Password   string  `json:"password" form:"password" binding:"required"`
	FirstName  string  `json:"first_name" form:"first_name" binding:"required"`
	LastName   string  `json:"last_name" form:"last_name" binding:"required"`
	Contact    string  `json:"contact" form:"contact" binding:"required"`
	Address    string  `json:"address" form:"address" binding:"required"`
	IsVerified bool    `json:"is_verified"`
	VerifiedAt *string `json:"verified_at"`
	CreatedAt  string  `json:"created_at"`
}

type OTPVerifyToken struct {
	UserID string `json:"id" form:"id" binding:"required"`
	OTP    string `json:"otp" form:"otp" binding:"required"`
}
