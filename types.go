package main

import "time"

type ROLE string

const (
	USER ROLE = "USER"
	ADMIN ROLE = "ADMIN"
)


type User struct {
	user_id int
	name string
	email string
	password string
	role ROLE 	
}

type Post struct {
	post_id int
	title string
	body string
	slug string
	author_id int
	created_at time.Time
	updated_at time.Time
}