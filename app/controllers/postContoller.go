package controllers

import (
	"encoding/json"
	"fmt"
	"go-news/consts"
	"go-news/database"
	"go-news/lib"
	"go-news/types"
	"go-news/views"
	"net/http"

	"github.com/gorilla/mux"
)
 




func CreatePost(w http.ResponseWriter, r * http.Request) {

	user, _ := r.Context().Value(consts.AUTH_COOKIE).(*types.User)
	
	post := types.Post {
		Title:r.FormValue("title"),
		Body:r.FormValue("body"),
		AuthorId:user.UserId,
		Slug:lib.Slugify(r.FormValue("title")),
	}


	if post.Title == ""{
		views.PostFormTmpl.Execute(w, struct {
			AuthData types.AuthData
			Post types.Post
			IsEditPage bool
			Success string
			Fail string
			TitleMessage string
			BodyMessage string
			}{
				AuthData:lib.FetchAuthData(r),
				Post:types.Post{},
				IsEditPage:false,
				Success:"",
				Fail:"",
				TitleMessage:"Title is required!",
				BodyMessage:"",
	
			})
	
		return
	}

	if  post.Body == ""  {
		views.PostFormTmpl.Execute(w, struct {
			AuthData types.AuthData
			Post types.Post
			IsEditPage bool
			Success string
			Fail string
			TitleMessage string
			BodyMessage string
			}{
				AuthData:lib.FetchAuthData(r),
				Post:types.Post{},
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
	_, err := database.GetDB().Exec(createPostQuery, post.Title, post.Body, post.AuthorId, post.Slug)

	if err != nil {
		views.PostFormTmpl.Execute(w, struct {
			AuthData types.AuthData
			Post types.Post
			IsEditPage bool
			Success string
			Fail string
			TitleMessage string
			BodyMessage string
			}{
				AuthData:lib.FetchAuthData(r),
				Post:types.Post{},
				IsEditPage:false,
				Success:"",
				Fail:"Unable to create,plz try again!",
				TitleMessage:"",
				BodyMessage:"",
			})
			return
	}

	// createPostId, _ := createPost.LastInsertId()
	// post.PostId = createPostId
	// json.NewEncoder(w).Encode(post)

	http.Redirect(w,r,"/me",http.StatusFound)
}

func PostEditPage(w http.ResponseWriter, r *http.Request) {
	post_id := mux.Vars(r)["post_id"]
	author, _ := r.Context().Value(consts.AUTH_COOKIE).(*types.User)

	var post types.Post

	isUserAdmin :=lib.IsAdmin(r)

	if isUserAdmin {
		err := database.GetDB().QueryRow(`
		select post_id, title, body, author_id
		from posts
		where post_id=?
	`,post_id).Scan(&post.PostId, &post.Title, &post.Body, &post.AuthorId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	} else {
		err := database.GetDB().QueryRow(`
		select post_id, title, body, author_id
		from posts
		where post_id=? and author_id=?
	`,post_id,author.UserId).Scan(&post.PostId, &post.Title, &post.Body, &post.AuthorId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	}
	




	views.PostFormTmpl.Execute(w, struct {
		AuthData types.AuthData
		Post types.Post
		IsEditPage bool
		Success string
		Fail string
		TitleMessage string
		BodyMessage string
		}{
			AuthData:lib.FetchAuthData(r),
			Post:post,
			IsEditPage:true,
			Success:"",
			Fail:"",
			TitleMessage:"",
			BodyMessage:"",
		})
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(consts.AUTH_COOKIE).(*types.User)
	user_id := user.UserId
	title := r.FormValue("title")
	body := r.FormValue("body")
	post_id := mux.Vars(r)["post_id"]

	if title == "" {
		views.PostFormTmpl.Execute(w, struct {
			AuthData types.AuthData
			Post types.Post
			IsEditPage bool
			Success string
			Fail string
			TitleMessage string
			BodyMessage string
			}{
				AuthData:lib.FetchAuthData(r),
				Post:types.Post{Title:title, Body:body,},
				IsEditPage:true,
				Success:"",
				Fail:"",
				TitleMessage:"Title is required!",
				BodyMessage:"",
			})
		return
	}

	if  body == ""  {
		views.PostFormTmpl.Execute(w, struct {
			AuthData types.AuthData
			Post types.Post
			IsEditPage bool
			Success string
			Fail string
			TitleMessage string
			BodyMessage string
			} {
				AuthData:lib.FetchAuthData(r),
				Post:types.Post{Title:title,Body:body,},
				IsEditPage:true,
				Success:"",
				Fail:"",
				TitleMessage:"",
				BodyMessage:"Body is required!",
			})
			return
	}
	selectStm := "select post_id from posts where author_id=? and post_id=?";

	isUserAdmin := lib.IsAdmin(r)
	var post_author_id int64
	if isUserAdmin {
		selectStm = "select post_id, author_id from posts where post_id=?";
		err  := database.GetDB().QueryRow(selectStm,post_id).Scan(&post_id, &post_author_id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}		
	} else {
		err  := database.GetDB().QueryRow(selectStm,user_id,post_id).Scan(&post_id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}		
	}
 
	updateStm := `update posts set title = ?, body = ? where post_id=?`

	_, err := database.GetDB().Exec(updateStm, title, body, post_id)

	if err != nil {
			views.PostFormTmpl.Execute(w, struct {
				AuthData types.AuthData
				Post types.Post
				IsEditPage bool
				Success string
				Fail string
				TitleMessage string
				BodyMessage string
				}{
					AuthData:lib.FetchAuthData(r),
					Post:types.Post{},
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

func DeletePost(w http.ResponseWriter, r * http.Request) {
	author, _ :=  r.Context().Value(consts.AUTH_COOKIE).(*types.User)
	author_id := author.UserId

	post_id := mux.Vars(r)["post_id"]

	deleteStm := `delete from posts where post_id=? and author_id=?`
 
	_, err := database.GetDB().Exec(deleteStm, post_id, author_id)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func DeletePostByAdmin(w http.ResponseWriter, r * http.Request) {

	post_id := mux.Vars(r)["post_id"]

	deleteStm := `delete from posts where post_id=?`

	_, err := database.GetDB().Exec(deleteStm, post_id)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func PostCreatePage(w http.ResponseWriter, r *http.Request) {
	views.PostFormTmpl.Execute(w, struct {
		AuthData types.AuthData
		Post types.Post
		IsEditPage bool
		Success string
		Fail string
		TitleMessage string
		BodyMessage string
		}{
			AuthData:lib.FetchAuthData(r),
			Post:types.Post{},
			IsEditPage:false,
			Success:"",
			Fail:"",
			TitleMessage:"",
			BodyMessage:"",

		})
}

func GetPostById(w http.ResponseWriter, r * http.Request) {
	post_id := mux.Vars(r)["post_id"]

	selectStm := `select p.post_id, p.title, p.slug, p.body, p.author_id, u.name
		from posts p
		inner join users u
		on p.author_id=u.user_id
		where p.post_id=?
	`;
	
	var post types.Post;
	err := database.GetDB().QueryRow(selectStm, post_id).Scan(&post.PostId, &post.Title, &post.Slug, &post.Body, &post.AuthorId, &post.Author)
	
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func GetPosts(w http.ResponseWriter, r * http.Request) {
	json.NewEncoder(w).Encode(lib.GetAllPostsUtil())
}

