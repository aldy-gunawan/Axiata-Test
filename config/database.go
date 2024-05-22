package config

import (
	"os"
	"fmt"
	"database/sql"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func DatabaseNew () *Database {
		db, err := sql.Open(os.Getenv("DB_TYPE"), os.Getenv("DB_DSN"))
		if err != nil {
			fmt.Println("Error opening database connection:", err)
			return nil
		}
	
		fmt.Println("Connected to the database!")
		return &Database{
			DB: db,
		}
}

func (d *Database) GetActiveDB() *sql.DB {
	return d.DB
}