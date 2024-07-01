package db

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// подключение к базе данных
var DB *sql.DB

func InitDB() error {
	var err error
	connStr := "user=postgres password=Fyfcnfcbz11 dbname=statsDB host=localhost port=5432 sslmode=disable"

	DB, err = sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("cant connect to the database: %v", err)
	}
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("cant ping the database: %v", err)
	}

	return nil
}
func CloseDB() {
	DB.Close()
}
