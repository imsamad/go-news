package main

import (
	"fmt"
	"go-news/controllers"
	"go-news/database"
	"go-news/lib"
	"go-news/middlewares"
	"go-news/types"
	"go-news/views"

	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	database.InitDB()

	r := mux.NewRouter()
	r.HandleFunc("/", homePage)

	r.HandleFunc("/me", middlewares.Chain(controllers.UserPanelPage, middlewares.AuthMiddleware()))

	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signup", controllers.Signup).Methods("POST")
	authRouter.HandleFunc("/signup", middlewares.Chain(controllers.SignupPage, middlewares.GuestMiddleware())).Methods("GET")
	authRouter.HandleFunc("/login", controllers.HandleLogin).Methods("POST")
	authRouter.HandleFunc("/login", middlewares.Chain(controllers.LoginPage, middlewares.GuestMiddleware())).Methods("GET")
	authRouter.HandleFunc("/me", middlewares.Chain(controllers.Me, middlewares.AuthMiddleware())).Methods("GET")
	authRouter.HandleFunc("/logout", controllers.Logout).Methods("GET")

	// TODO: figure out better route names conventions
	postRouter := r.PathPrefix("/posts").Subrouter()
	postRouter.HandleFunc("", middlewares.Chain(controllers.PostCreatePage, middlewares.AuthMiddleware())).Methods("GET")
	postRouter.HandleFunc("", middlewares.Chain(controllers.CreatePost, middlewares.AuthMiddleware())).Methods("POST")
	postRouter.HandleFunc("/all", controllers.GetPosts).Methods("GET")
	// because golang does not have redirect with ctx
	// have to go with this bad convention for the moment
	postRouter.HandleFunc("/edit/{post_id}", middlewares.Chain(controllers.PostEditPage, middlewares.AuthMiddleware())).Methods("GET")
	postRouter.HandleFunc("/edit/{post_id}", middlewares.Chain(controllers.UpdatePost, middlewares.AuthMiddleware())).Methods("POST")

	postRouter.HandleFunc("/{post_id}", middlewares.Chain(controllers.DeletePost, middlewares.AuthMiddleware())).Methods("DELETE")
	postRouter.HandleFunc("/admin/{post_id}", middlewares.Chain(controllers.DeletePostByAdmin, middlewares.AdminMiddleware())).Methods("DELETE")
	postRouter.HandleFunc("/single/{post_id}", controllers.GetPostById).Methods("GET")

	r.HandleFunc("/admin/", middlewares.Chain(controllers.GetAdminPage, middlewares.AdminMiddleware())).Methods("GET")

	fmt.Println("app is running on 3000")

	http.ListenAndServe(":3000", r)
}

func homePage(w http.ResponseWriter, r *http.Request) {

	type Posts struct {
		Posts    []types.Post
		AuthData types.AuthData
	}

	posts := Posts{
		Posts:    lib.GetAllPostsUtil(),
		AuthData: lib.FetchAuthData(r),
	}

	views.HomePageTmpl.Execute(w, posts)
}
