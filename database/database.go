package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	fmt.Println(connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("here")
		return nil, err
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		fmt.Println("or here")

		return nil, err
	}

	// Set connection pool settings (optional tapi recommended)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connected successfully")
	return db, nil
}
