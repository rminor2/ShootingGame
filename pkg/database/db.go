package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectToDatabase() (*sql.DB, error) {
	// Replace with your actual database credentials
	const (
		host     = "your-database-host.amazonaws.com"
		port     = 5432
		user     = "yourDatabaseUser"
		password = "yourDatabasePassword"
		dbname   = "yourDatabaseName"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
