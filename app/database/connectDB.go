package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	godotenv.Load("../.env")

	db, err := sql.Open("mysql", os.Getenv("MY_SQL_URI"))

	if err != nil {
		fmt.Println("error mysql conn..")
		fmt.Println("reason: ", err)
		os.Exit(1)
	}

	DB = db
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return DB
}
