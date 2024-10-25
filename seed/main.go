package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	godotenv.Load("../.env")
	db, err := sql.Open("mysql", os.Getenv("MY_SQL_URI"))

	if err != nil {
		fmt.Println("unable to connect to mysql")
		fmt.Println("here is the reason: ", err)
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		fmt.Println("error on ping : ", err)
		os.Exit(1)
	}
	// Create the database if it doesn't exist
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS news")
	if err != nil {
		fmt.Println("Error creating database")
		fmt.Println(err)
		os.Exit(1)
	}
	userTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
		user_id INT AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		role ENUM("USER", "ADMIN") NOT NULL DEFAULT "USER",
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY(user_id)		
		);
	`

	if _, err := db.Exec(userTableQuery); err != nil {
		fmt.Println("error creating user table")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("users table created")

	postTableQuery := `
		CREATE TABLE IF NOT EXISTS posts (
			post_id INT AUTO_INCREMENT,
			title VARCHAR(255) NOT NULL,
			slug VARCHAR(255) NOT NULL UNIQUE,
			body VARCHAR(255) NOT NULL,
			author_id INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY(post_id),
			FOREIGN KEY(author_id) REFERENCES users(user_id)
		);
	`

	if _, err = db.Exec(postTableQuery); err != nil {
		fmt.Println("error creating user table")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("posts table created")

	// User Seed
	{
		passwordHashedBytes, err := bcrypt.GenerateFromPassword([]byte("123456"), 14)

		if err != nil {
			fmt.Println("unable to hash pwd, that wht exiting out...")
			os.Exit(1)
		}

		password := string(passwordHashedBytes)

		users := []struct {
			name     string
			email    string
			password string
			role     string
		}{
			{email: "admin@email.com", password: password, name: "samad", role: "ADMIN"},
			{email: "user1@email.com", password: password, name: "samad1", role: "USER"},
			{email: "user2@email.com", password: password, name: "samad2", role: "USER"},
			{email: "user3@email.com", password: password, name: "samad3", role: "USER"},
			{email: "user4@email.com", password: password, name: "samad4", role: "USER"},
		}

		for _, user := range users {
			_, err := db.Exec("insert into users (name, email,password,role) values (?,?,?,?)", user.name, user.email, user.password, user.role)
			if err != nil {
				fmt.Println("error inserting row of :", user.email)
				fmt.Println(err)
			}
		}

		fmt.Println("5 users inserted successfully")
	}

	// Posts seed
	{
		rows, err := db.Query("select user_id from users")
		defer rows.Close()

		if err != nil {
			fmt.Println("error while fetching user ids")
			os.Exit(1)
		}

		var users []int

		for rows.Next() {
			var id int

			if err := rows.Scan(&id); err != nil {
				fmt.Println("error while scanning user rows")
			}
			users = append(users, id)
		}

		for id := range users {
			for i := 0; i < 10; i++ {
				post := struct {
					title     string
					body      string
					slug      string
					author_id int
				}{
					title:     fmt.Sprintf("User %d post - %d", id, i),
					slug:      fmt.Sprintf("user-%d-post-%d", id, i),
					body:      fmt.Sprintf("This is the body of post by User %d.", id),
					author_id: id,
				}

				_, err := db.Exec("insert into posts (title, body, slug, author_id) values (?,?,?,?)", post.title, post.body, post.slug, post.author_id)

				if err != nil {
					fmt.Println("error while seeding post")
					fmt.Println(err)
				}

				fmt.Println(post)

			}
		}

		fmt.Println("users and post seeded")

	}
}
