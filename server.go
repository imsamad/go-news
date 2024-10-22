package main

import (
	"github.com/gorilla/mux"
	"net/http"

    "database/sql"
	_ "github.com/go-sql-driver/mysql"

	"os"
	"fmt"

	"github.com/gorilla/sessions"
)

var DB *sql.DB

var (
	key = []byte("super-secret")
	store = sessions.NewCookieStore(key)
)
const AUTH_COOKIE string ="AUTH_COOKIE"
func main() {

	r := mux.NewRouter()
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signup",signup).Methods("POST")
	authRouter.HandleFunc("/login",handleLogin).Methods("POST")
	authRouter.HandleFunc("/me",me).Methods("GET")
	authRouter.HandleFunc("/logout",logout).Methods("GET")

	db, err := sql.Open("mysql", "root:root@(127.0.0.1:3306)/news1")

	if err != nil {
		fmt.Println("error mysql conn..")
		os.Exit(1)
	}
	DB = db

	http.ListenAndServe(":3000",r)
}