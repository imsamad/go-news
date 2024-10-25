package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	db, err := sql.Open("mysql", "root:root@(127.0.0.1:3306)/news4")

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
