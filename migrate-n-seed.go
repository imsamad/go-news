package main

import (
	"fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
	"os"
)

func main() {
	db,err := sql.Open("mysql","root:root@(127.0.0.1:3306)/news1")

	if err != nil {
		fmt.Println("unable to connect to mysql")
		fmt.Println("here is the reason: ",err)
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		fmt.Println("error on ping")
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
		password,_ := HashPassword("123456")
		users := []User {
			{ email:"admin@gmail.com",password:password,name:"samad" },
			{ email:"user1@gmail.com",password:password,name:"samad1" },
			{ email:"user2@gmail.com",password:password,name:"samad2" },
			{ email:"user3@gmail.com",password:password,name:"samad3" },
			{ email:"user4@gmail.com",password:password,name:"samad4" },
		}

		for _,user := range users {
			_, err := db.Exec("insert into users (name, email,password) values (?,?,?)",user.name, user.email, user.password)
			if err != nil {
				fmt.Println("error inserting row of :",user.email)
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
			post := Post {
				title: fmt.Sprintf("User %d post",id),
				slug: fmt.Sprintf("user-%d-post",id),
				body : fmt.Sprintf("This is the body of post by User %d.", id),
				author_id: id,
			}

			_, err := db.Exec("insert into posts (title, body, slug, author_id) values (?,?,?,?)",post.title, post.body, post.slug, post.author_id)

			if err != nil {
				fmt.Println("error while seeding post")
				fmt.Println(err)
			}

			fmt.Println(post)
		}

		fmt.Println("users and post seeded")

	}
}