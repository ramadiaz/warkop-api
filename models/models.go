package models

import (
	"time"
)

type ClientTrack struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	IP        string    `gorm:"type:varchar(255)"`
	Browser   string    `gorm:"type:varchar(255)"`
	Version   string    `gorm:"type:varchar(255)"`
	OS        string    `gorm:"type:varchar(255)"`
	Device    string    `gorm:"type:varchar(255)"`
	Origin    string    `gorm:"type:varchar(255)"`
	API       string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (ClientTrack) TableName() string {
	return "client_track"
}

type User struct {
	ID         string     `gorm:"type:uuid;primaryKey;default:uuid_generate_v1()"`
	Username   string     `gorm:"type:varchar(255);unique;not null"`
	Email      string     `gorm:"type:varchar(255);not null"`
	Password   string     `gorm:"type:varchar(255);not null"`
	FirstName  string     `gorm:"type:varchar(255);not null"`
	LastName   string     `gorm:"type:varchar(255);not null"`
	Contact    string     `gorm:"type:varchar(255);not null"`
	Address    string     `gorm:"type:varchar(255);not null"`
	IsVerified bool       `gorm:"default:false"`
	VerifiedAt *time.Time `gorm:"type:timestamp"`
	CreatedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
}

func (User) TableName() string {
	return "users"
}

type VerificationToken struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	UserID    string    `gorm:"type:uuid;not null"`
	Token     string    `gorm:"type:varchar(255);not null"`
	ExpiredAt time.Time `gorm:"type:timestamp;not null"`
}

func (VerificationToken) TableName() string {
	return "verification_token"
}

type ResetOTP struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	UserID    string    `gorm:"type:uuid;not null"`
	OTP       int       `gorm:"not null"`
	ExpiredAt time.Time `gorm:"type:timestamp;not null"`
}

func (ResetOTP) TableName() string {
	return "reset_otp"
}

type APIKey struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Token     string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (APIKey) TableName() string {
	return "api_key"
}

type Menu struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Type      string    `gorm:"type:menu_type;not null"`
	Price     int       `gorm:"not null;default:0"`
	Stock     int       `gorm:"not null;default:0"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (Menu) TableName() string {
	return "menu"
}

type Transaction struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	CashierID string    `gorm:"type:varchar(255);not null"`
	Total     int64     `gorm:"not null"`
	Cash      int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (Transaction) TableName() string {
	return "transaction"
}

type TransactionItem struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement"`
	TransactionID int64     `gorm:"not null"`
	MenuID        int       `gorm:"not null"`
	Quantity      int       `gorm:"not null"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (TransactionItem) TableName() string {
	return "transaction_item"
}

type UsersImage struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	UserID    string    `gorm:"type:varchar(255);unique;not null"`
	ImageURL  string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (UsersImage) TableName() string {
	return "users_image"
}
