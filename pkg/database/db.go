package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectToDatabase() (*sql.DB, error) {
	const (
		host     = "shootdb-1.czm0a6c2szh6.us-east-2.rds.amazonaws.com"
		port     = 5432
		user     = "postgres"
		password = "Dragon123"
		dbname   = "shootdb"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected!")
	return db, nil
}
