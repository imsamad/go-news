package main

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
)

var (
    postFormTmpl    = mustParseTemplate("views/post_form.tmpl")
)

func getAllPostsUtil() []Post {
	var posts []Post

	rows, err := DB.Query(`
		select p.post_id, p.title, p.slug, p.body, p.author_id, u.name
		from posts p
		inner join users u
		on p.author_id=u.user_id
	`)

	defer rows.Close()

	if err != nil {
		return posts
	}

	for rows.Next() {
		var post Post
		_ = rows.Scan(&post.PostId, &post.Title, &post.Slug, &post.Body, &post.AuthorId, &post.Author)

		posts = append(posts, post)
	}

	return posts
}

func getAllPostsExceptUtil(user_id int64) []Post {
	var posts []Post

	rows, err := DB.Query(`
		select p.post_id, p.title, p.slug, p.body, p.author_id, u.name
		from posts p
		inner join users u
		on p.author_id=u.user_id
		where not p.author_id=? 
	`,user_id)

	defer rows.Close()

	if err != nil {
		return posts
	}

	for rows.Next() {
		var post Post
		_ = rows.Scan(&post.PostId, &post.Title, &post.Slug, &post.Body, &post.AuthorId, &post.Author)

		posts = append(posts, post)
	}

	return posts
}

func getAllMyPostsUtil(user_id int64) []Post {
	var posts []Post

	rows, err := DB.Query(`
		select p.post_id, p.title, p.slug, p.body, p.author_id, u.name
		from posts p
		inner join users u
		on p.author_id=u.user_id
		where user_id=?
	`,user_id)

	defer rows.Close()

	if err != nil {
		return posts
	}

	for rows.Next() {
		var post Post
		_ = rows.Scan(&post.PostId, &post.Title, &post.Slug, &post.Body, &post.AuthorId, &post.Author)

		posts = append(posts, post)
	}

	return posts
}

func createPost(w http.ResponseWriter, r * http.Request) {

	user, _ := r.Context().Value(AUTH_COOKIE).(*User)
	
	post := Post {
		Title:r.FormValue("title"),
		Body:r.FormValue("body"),
		AuthorId:user.UserId,
		Slug:Slugify(r.FormValue("title")),
	}


	if post.Title == ""{
		postFormTmpl.Execute(w, struct {
			AuthData AuthData
			Post Post
			IsEditPage bool
			Success string
			Fail string
			TitleMessage string
			BodyMessage string
			}{
				AuthData:fetchAuthData(r),
				Post:Post{},
				IsEditPage:false,
				Success:"",
				Fail:"",
				TitleMessage:"Title is required!",
				BodyMessage:"",
	
			})
	
		return
	}

	if  post.Body == ""  {
		postFormTmpl.Execute(w, struct {
			AuthData AuthData
			Post Post
			IsEditPage bool
			Success string
			Fail string
			TitleMessage string
			BodyMessage string
			}{
				AuthData:fetchAuthData(r),
				Post:Post{},
				IsEditPage:false,
				Success:"",
				Fail:"",
				TitleMessage:"",
				BodyMessage:"Body is required!",
			})
			return
	}

	createPostQuery := `
	insert into posts (title,body,author_id,slug) values (?,?,?,?);
	`
	createPost, err := DB.Exec(createPostQuery, post.Title, post.Body, post.AuthorId, post.Slug)

	if err != nil {
		postFormTmpl.Execute(w, struct {
			AuthData AuthData
			Post Post
			IsEditPage bool
			Success string
			Fail string
			TitleMessage string
			BodyMessage string
			}{
				AuthData:fetchAuthData(r),
				Post:Post{},
				IsEditPage:false,
				Success:"",
				Fail:"Unable to create,plz try again!",
				TitleMessage:"",
				BodyMessage:"",
			})
			return
	}

	createPostId, _ := createPost.LastInsertId()
	post.PostId = createPostId
	// json.NewEncoder(w).Encode(post)

	http.Redirect(w,r,"/me",http.StatusFound)
}

func postEditPage(w http.ResponseWriter, r *http.Request) {
	post_id := mux.Vars(r)["post_id"]
	author, _ := r.Context().Value(AUTH_COOKIE).(*User)

	var post Post

	isUserAdmin := isAdmin(r)

	if isUserAdmin {
		err := DB.QueryRow(`
		select post_id, title, body, author_id
		from posts
		where post_id=?
	`,post_id).Scan(&post.PostId, &post.Title, &post.Body, &post.AuthorId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	} else {
		err := DB.QueryRow(`
		select post_id, title, body, author_id
		from posts
		where post_id=? and author_id=?
	`,post_id,author.UserId).Scan(&post.PostId, &post.Title, &post.Body, &post.AuthorId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	}
	




	postFormTmpl.Execute(w, struct {
		AuthData AuthData
		Post Post
		IsEditPage bool
		Success string
		Fail string
		TitleMessage string
		BodyMessage string
		}{
			AuthData:fetchAuthData(r),
			Post:post,
			IsEditPage:true,
			Success:"",
			Fail:"",
			TitleMessage:"",
			BodyMessage:"",
		})
		return
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(AUTH_COOKIE).(*User)
	user_id := user.UserId
	title := r.FormValue("title")
	body := r.FormValue("body")
	post_id := mux.Vars(r)["post_id"]

	if title == "" {
		postFormTmpl.Execute(w, struct {
			AuthData AuthData
			Post Post
			IsEditPage bool
			Success string
			Fail string
			TitleMessage string
			BodyMessage string
			}{
				AuthData:fetchAuthData(r),
				Post:Post{Title:title, Body:body,},
				IsEditPage:true,
				Success:"",
				Fail:"",
				TitleMessage:"Title is required!",
				BodyMessage:"",
			})
		return
	}

	if  body == ""  {
		postFormTmpl.Execute(w, struct {
			AuthData AuthData
			Post Post
			IsEditPage bool
			Success string
			Fail string
			TitleMessage string
			BodyMessage string
			} {
				AuthData:fetchAuthData(r),
				Post:Post{Title:title,Body:body,},
				IsEditPage:true,
				Success:"",
				Fail:"",
				TitleMessage:"",
				BodyMessage:"Body is required!",
			})
			return
	}
	selectStm := "select post_id from posts where author_id=? and post_id=?";

	isUserAdmin := isAdmin(r)
	var post_author_id int64
	if isUserAdmin {
		selectStm = "select post_id, author_id from posts where post_id=?";
		err  := DB.QueryRow(selectStm,post_id).Scan(&post_id, &post_author_id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}		
	} else {
		err  := DB.QueryRow(selectStm,user_id,post_id).Scan(&post_id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}		
	}
 
	updateStm := `update posts set title = ?, body = ? where post_id=?`

	_, err := DB.Exec(updateStm, title, body, post_id)

	if err != nil {
			postFormTmpl.Execute(w, struct {
				AuthData AuthData
				Post Post
				IsEditPage bool
				Success string
				Fail string
				TitleMessage string
				BodyMessage string
				}{
					AuthData:fetchAuthData(r),
					Post:Post{},
					IsEditPage:false,
					Success:"",
					Fail:"Operation failed, please try again!",
					TitleMessage:"",
					BodyMessage:"",
				})
				return
	}

	// check if it is admin and post is being changed not belongs to admin
	// then it quite evident update post request is being orignited from admin panel
	// so redirect the req flow to admin panel instead of user panel
	fmt.Println("post_author_id: ",post_author_id)
	fmt.Println("user_id: ",user_id)

	if post_author_id != 0 && post_author_id != user_id {
		fmt.Println("redirecting to admin panel...")
		http.Redirect(w,r,"/admin/",http.StatusFound)
	}else {
		http.Redirect(w,r,"/me",http.StatusFound)
	}
}

func deletePost(w http.ResponseWriter, r * http.Request) {
	author, _ :=  r.Context().Value(AUTH_COOKIE).(*User)
	author_id := author.UserId

	post_id := mux.Vars(r)["post_id"]

	deleteStm := `delete from posts where post_id=? and author_id=?`
 
	_, err := DB.Exec(deleteStm, post_id, author_id)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func deletePostByAdmin(w http.ResponseWriter, r * http.Request) {

	post_id := mux.Vars(r)["post_id"]

	deleteStm := `delete from posts where post_id=?`

	_, err := DB.Exec(deleteStm, post_id)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func postCreatePage(w http.ResponseWriter, r *http.Request) {
	postFormTmpl.Execute(w, struct {
		AuthData AuthData
		Post Post
		IsEditPage bool
		Success string
		Fail string
		TitleMessage string
		BodyMessage string
		}{
			AuthData:fetchAuthData(r),
			Post:Post{},
			IsEditPage:false,
			Success:"",
			Fail:"",
			TitleMessage:"",
			BodyMessage:"",

		})
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
	json.NewEncoder(w).Encode(getAllPostsUtil())
}

