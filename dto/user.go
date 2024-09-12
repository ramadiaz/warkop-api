package dto

type User struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Contact    string `json:"contact"`
	Address    string `json:"address"`
	IsVerified bool   `json:"is_verified"`
	VerifiedAt string `json:"verified_at"`
	CreatedAt  string `json:"created_at"`
}
