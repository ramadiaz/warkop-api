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

	return &compRepository{
		DB: db,
	}
}
