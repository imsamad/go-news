package types

import "time"

type ROLE string
const (
	USER ROLE = "USER"
	ADMIN ROLE = "ADMIN"
)

type User struct {
	UserId int64 `json:"user_id"`
	Name string `json:"name"`
	Email string `json:"email"`
	password string 
	Role ROLE 	`json:"role"`
}

type Post struct {
	PostId int64 `json:"post_id"`
	Title string `json:"title"`
	Body string `json:"body"`
	Slug string `json:"slug"`
	AuthorId int64 `json:"author_id"`	
	Author string `json:author`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthData struct {
	IsLoggedIn bool
	IsAdmin bool
}