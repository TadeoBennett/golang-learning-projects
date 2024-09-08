package initializers

import (
	"database/sql"
	"log"
)

func ConnectToDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Println("Ping Successful")
		return nil, err
	}
	log.Println("Database Connected")
	return db, nil
}
