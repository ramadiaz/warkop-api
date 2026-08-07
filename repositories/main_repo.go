package repositories

import (
	"log"
	"warkop-api/config"
	"warkop-api/dto"
	"warkop-api/models"

	"gorm.io/gorm"
)

type CompRepository interface {
	RegisterUser(data dto.User) (*string, error)
	RegisterToken(data dto.User) (*string, error)
	RegisterAPIKey(name string, secret string) error
	VerifyAccount(token string) error
	GetUserData(username string) (*dto.User, error)
	UploadUserProfile(data dto.User, image_url string) error
	GetUserProfile(id string) (*string, error)

	RegisterMenu(data dto.Menu) error
	GetAllMenu() ([]*dto.Menu, error)

	RegisterTransaction(data dto.Transaction) (*int64, error)
	RegisterTransactionItem(data dto.TransactionItem) error
	GetTransaction(id string) (*dto.Transaction, error)
	GetTransactionItemInTx(tx *gorm.DB, id string) ([]*dto.TransactionItem, error)
	GetTransactionItem(id string) ([]*dto.TransactionItem, error)
	GetAllTransaction() ([]*dto.Transaction, error)

	RequestResetPassword(data dto.User, otp string) error
	VerifyResetPassword(data dto.OTPVerifyToken) (*dto.OTPVerifyToken, error)
	ResetPassword(user_data dto.User) error

	BeginTransaction() (*gorm.DB, error)
}

type compRepository struct {
	DB *gorm.DB
}

func NewComponentRepository(DB *gorm.DB) *compRepository {
	db := config.InitDB()

	// Enable uuid-ossp extension
	err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error
	if err != nil {
		log.Fatalf("Error creating extension: %v", err)
	}

	// Create custom enum type menu_type if not exists
	err = db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'menu_type') THEN
				CREATE TYPE menu_type AS ENUM ('Food', 'Drink', 'Snack', 'Other');
			END IF;
		END $$;
	`).Error
	if err != nil {
		log.Fatalf("Error creating enum type: %v", err)
	}

	// Auto-migrate models
	err = db.AutoMigrate(
		&models.ClientTrack{},
		&models.User{},
		&models.VerificationToken{},
		&models.ResetOTP{},
		&models.APIKey{},
		&models.Menu{},
		&models.Transaction{},
		&models.TransactionItem{},
		&models.UsersImage{},
	)
	if err != nil {
		log.Fatalf("Error running auto-migration: %v", err)
	}

	return &compRepository{
		DB: db,
	}
}
