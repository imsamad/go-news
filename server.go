package main

import (
	"github.com/gorilla/mux"
	"net/http"

    "database/sql"
	_ "github.com/go-sql-driver/mysql"

	"os"
	"fmt"

	"github.com/gorilla/sessions"
	
	"html/template"
)

var DB *sql.DB

var (
	key = []byte("super-secret")
	store = sessions.NewCookieStore(key)
)

const AUTH_COOKIE string ="AUTH_COOKIE"

func main() {
	// homepage
	// signup page
	// login page
	// user panel posts table
	// delete
	// create story page
	// edit story page

	// admin panel

	r := mux.NewRouter()

	homePageTmpl, err :=  template.ParseFiles("views/index.tmpl")

	if err != nil {
		fmt.Println(err)
	}

	r.HandleFunc("/",func(w http.ResponseWriter, r *http.Request) {	
		
		type Posts struct {
			Posts []Post
			AuthData AuthData
		}	

		posts := Posts{
			Posts: getAllPostsUtil(),
			AuthData:fetchAuthData(r),
		}
		
		err = homePageTmpl.Execute(w,posts)
	})
	
	r.HandleFunc("/me",middlewares(userPanelPage,authMiddleware()))
	
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signup",signup).Methods("POST")
	authRouter.HandleFunc("/signup",middlewares(signupPage,guestMiddleware())).Methods("GET")
	authRouter.HandleFunc("/login",handleLogin).Methods("POST")
	authRouter.HandleFunc("/login",middlewares(loginPage,guestMiddleware())).Methods("GET")
	authRouter.HandleFunc("/me",middlewares(me, authMiddleware())).Methods("GET")
	authRouter.HandleFunc("/logout",logout).Methods("GET")

	// TODO: figure out better route names conventions
	postRouter := r.PathPrefix("/posts").Subrouter()
	postRouter.HandleFunc("",middlewares(postCreatePage, authMiddleware())).Methods("GET")
	postRouter.HandleFunc("",middlewares(createPost, authMiddleware())).Methods("POST")

	postRouter.HandleFunc("/all",getPosts).Methods("GET")
	
	// because golang does not have redirect with ctx
	// have to go with this bad convention for the moment
	postRouter.HandleFunc("/edit/{post_id}",middlewares(postEditPage, authMiddleware())).Methods("GET")
	postRouter.HandleFunc("/edit/{post_id}",middlewares(updatePost, authMiddleware())).Methods("POST")
	
	postRouter.HandleFunc("/{post_id}",middlewares(deletePost, authMiddleware())).Methods("DELETE")
	postRouter.HandleFunc("/admin/{post_id}",middlewares(deletePostByAdmin, adminMiddleware())).Methods("DELETE")
	postRouter.HandleFunc("/single/{post_id}",getPostById).Methods("GET")

	r.HandleFunc("/admin/",middlewares(getAdminPage, adminMiddleware())).Methods("GET")

	db, err := sql.Open("mysql", "root:root@(127.0.0.1:3306)/news4")

	if err != nil {
		fmt.Println("error mysql conn..")
		fmt.Println("reason: ",err)
		os.Exit(1)
	}
	DB = db

	http.ListenAndServe(":3000",r)
}