package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
)

func createPost(w http.ResponseWriter, r * http.Request) {

	user, _ := r.Context().Value(AUTH_COOKIE).(*User)
	
	post := Post {
		Title:r.FormValue("title"),
		Body:r.FormValue("body"),
		AuthorId:user.UserId,
		Slug:Slugify(r.FormValue("title")),
	}
	if post.Title == "" || post.Body == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return
	}
	createPostQuery := `
	insert into posts (title,body,author_id,slug) values (?,?,?,?);
	`

	createPost, err := DB.Exec(createPostQuery, post.Title, post.Body, post.AuthorId, post.Slug)

	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	createPostId, _ := createPost.LastInsertId()
	post.PostId = createPostId
	json.NewEncoder(w).Encode(post)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	author, _ := r.Context().Value(AUTH_COOKIE).(*User)
	author_id := author.UserId
	title := r.FormValue("title")
	body := r.FormValue("body")
	post_id := mux.Vars(r)["post_id"]

	if title == "" || body == "" {
		fmt.Println("author_id ",author_id)
		http.Error(w, http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return
	}

	err  := DB.QueryRow("select post_id from posts where author_id=? and post_id=?",author_id,post_id).Scan(&post_id)
	fmt.Println("one",err)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	updateStm := `update posts set title = ?, body = ? where post_id=?`

	_, err = DB.Exec(updateStm, title, body, post_id)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func deletePost(w http.ResponseWriter, r * http.Request) {
	author, _ :=  r.Context().Value(AUTH_COOKIE).(*User)
	author_id := author.UserId

	post_id := mux.Vars(r)["post_id"]

	deleteStm := `delete from posts where post_id=? and author_id=?`

	result, err := DB.Exec(deleteStm, post_id, author_id)

	fmt.Println(result)
	fmt.Println(err)
}

func getPostById(w http.ResponseWriter, r * http.Request) {
	post_id := mux.Vars(r)["post_id"]

	selectStm := `select p.post_id, p.title, p.slug, p.body, p.author_id, u.name
		from posts p
		inner join users u
		on p.author_id=u.user_id
		where p.post_id=?
	`;
	
	var post Post;
	err := DB.QueryRow(selectStm, post_id).Scan(&post.PostId, &post.Title, &post.Slug, &post.Body, &post.AuthorId, &post.Author)
	
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func getPosts(w http.ResponseWriter, r * http.Request) {

	rows, err := DB.Query(`select p.post_id, p.title, p.slug, p.body, p.author_id, u.name
		from posts p
		inner join users u
		on p.author_id=u.user_id
	`)

	defer rows.Close()

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var posts []Post

	for rows.Next() {
		var post Post
		_ = rows.Scan(&post.PostId, &post.Title, &post.Slug, &post.Body, &post.AuthorId, &post.Author)

		posts = append(posts, post)
	}

	json.NewEncoder(w).Encode(posts)
}