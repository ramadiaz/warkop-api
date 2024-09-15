package repositories

import (
	"database/sql"
	"log"
	"warkop-api/config"
	"warkop-api/dto"
)

type CompRepository interface {
	RegisterUser(data dto.User) (*string, error)
	RegisterToken(data dto.User) (*string, error)
	RegisterAPIKey(name string, secret string) error
	VerifyAccount(token string) error
	GetUserData(username string) (*dto.User, error)

	RegisterMenu(data dto.Menu) error
	GetAllMenu() ([]*dto.Menu, error)

	RegisterTransaction(data dto.Transaction) (*int64, error)
	RegisterTransactionItem(data dto.TransactionItem) error 
	GetTransaction(id string) (*dto.Transaction, error) 
	GetTransactionItem(id string) ([]*dto.TransactionItem, error)
}

type compRepository struct {
	DB *sql.DB
}

func NewComponentRepository(DB *sql.DB) *compRepository {
	db := config.InitDB()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS client_track (
			id 			BIGSERIAL PRIMARY KEY NOT NULL,
			ip 			VARCHAR(255),
			browser 	VARCHAR(255),
			version 	VARCHAR(255),
			os 			VARCHAR(255),
			device 		VARCHAR(255),
			origin 		VARCHAR(255),
			api 		VARCHAR(255),
			created_at 	TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	_, err = db.Exec(`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

		CREATE TABLE IF NOT EXISTS users (
			id          UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v1(),
			username 	VARCHAR(255) UNIQUE NOT NULL,
			email 		VARCHAR(255) NOT NULL,
			password 	VARCHAR(255) NOT NULL,
			first_name 	VARCHAR(255) NOT NULL,
			last_name 	VARCHAR(255) NOT NULL,
			contact 	VARCHAR(255) NOT NULL,
			address 	VARCHAR(255) NOT NULL,
			is_verified BOOLEAN DEFAULT FALSE,
			verified_at TIMESTAMP,
			created_at 	TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS verification_token (
			id 			BIGSERIAL PRIMARY KEY NOT NULL,
			user_id 	UUID NOT NULL,
			token 		VARCHAR(255) NOT NULL,
			expired_at  TIMESTAMP NOT NULL
		);`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS api_key (
			id 			BIGSERIAL PRIMARY KEY NOT NULL,
			name 		VARCHAR(255) NOT NULL,
			token 		VARCHAR(255) NOT NULL,
			created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	_, err = db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'menu_type') THEN
				CREATE TYPE menu_type AS ENUM ('Food', 'Drink', 'Snack', 'Other');
			END IF;
		END $$;

		CREATE TABLE IF NOT EXISTS menu (
			id          BIGSERIAL PRIMARY KEY NOT NULL,
			name        VARCHAR(255) NOT NULL,
			type        menu_type NOT NULL,
			price       INT NOT NULL DEFAULT 0,
			stock       INT NOT NULL DEFAULT 0,
			created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS transaction (
			id 			BIGSERIAL PRIMARY KEY NOT NULL,
			cashier_id 	VARCHAR(255) NOT NULL,
			total 		INT NOT NULL,
			cash        INT NOT NULL,
			created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS transaction_item (
			id 				BIGSERIAL PRIMARY KEY NOT NULL,
			transaction_id 	INT NOT NULL,
			menu_id 		INT NOT NULL,
			quantity 		INT NOT NULL,
			created_at  	TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	return &compRepository{
		DB: db,
	}
}
