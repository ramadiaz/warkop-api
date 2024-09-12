package repositories

import (
	"database/sql"
	"log"
	"warkop-api/config"
	"warkop-api/dto"
)

type CompRepository interface {
	RegisterUser(data dto.User) (int64, error)
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
		CREATE TABLE IF NOT EXISTS user (
			id 			BIGSERIAL PRIMARY KEY NOT NULL,
			email 		VARCHAR(255) UNIQUE,
			username 	VARCHAR(255) UNIQUE,
			password 	VARCHAR(255),
			first_name 	VARCHAR(255),
			last_name 	VARCHAR(255),
			contact 	VARCHAR(255),
			address 	VARCHAR(255),
			is_verified BOOLEAN DEFAULT FALSE,
			verified_at TIMESTAMP,
			created_at 	TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	return &compRepository{
		DB: db,
	}
}
