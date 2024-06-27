package db

import (
    "database/sql"
    "fmt"
    _ "github.com/jackc/pgx/v4/stdlib"
)
var DB *sql.DB
func InitDB()error{
	var err error
	connStr:= fmt.Sprintf(
		"user=%s  host=%s port=%d dbname=%s sslmode=disable",
		"postgres", "localhost", 5432, "statsDB",
	  )
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