package main

import (
	"net/http"
	"fmt"
	"encoding/json"
)

type LoginUser struct {
	UserId int64 `json:"user_id"`
	Name string `json:"name"`
	Email string `json:"email"`
	password string 
}

func signup(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	name := r.FormValue("name")
	password := r.FormValue("password")
	
	if email == "" || password == "" || name == "" {
		fmt.Fprintln(w, http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return
	} 

	alreadyExistedQuery := `
		select email from users where email=?
	`

	notExist := DB.QueryRow(alreadyExistedQuery, email).Scan(&email)
	// fmt.Println("alreadyExistedUser: ",alreadyExistedUser)	
	if notExist == nil {
		fmt.Fprintln(w, "Email already exist",404)
		return
	}	

	userQuery := `
		insert into users (name, email,password) values (?,?,?);
	`
	
    hashPwd, _ := HashPassword(password)		
	newUser, err := DB.Exec(userQuery, name, email, hashPwd)


	if err != nil {
		fmt.Println("err: ",err)
		fmt.Fprintln(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}	
	
	id, err := newUser.LastInsertId()
	
	if err != nil {
		fmt.Println("err fetching id of newly created user: ",err)
		fmt.Fprintln(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	
	session, _ := store.Get(r, AUTH_COOKIE)

	session.Values[AUTH_COOKIE] = id
	session.Save(r, w)
	 
	json.NewEncoder(w).Encode(LoginUser {Name:name, Email:email, UserId:id})
}

func handleLogin(w http.ResponseWriter, r * http.Request) {
	
	user := LoginUser{
		Email :  r.FormValue("email"),
		password :  r.FormValue("password"),
	}
	
	if user.Email == "" || user.password == "" {
		fmt.Fprintln(w, http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return
	} 
	
	var dbPwd string;
	var user_id int64;
	var name string;

	err := DB.QueryRow("select user_id, email, name, password from users where email=?",user.Email).Scan(&user_id, &user.Email, &name, &dbPwd)

	if err != nil {
		fmt.Println(err)
		return
	}

	if !CheckPasswordHash(user.password, dbPwd) {
		fmt.Println(err)
		return
	}

	user.UserId = user_id
	user.Name = name

	session , err := store.Get(r, AUTH_COOKIE)
 
	session.Values[AUTH_COOKIE] = user_id
	session.Save(r,w)
	json.NewEncoder(w).Encode(user)
}

func me(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, AUTH_COOKIE)
	user_id, ok := session.Values[AUTH_COOKIE].(int64)
	
	if  !ok || user_id == -1 {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
 
	user := LoginUser {
		UserId:user_id,
	}

	err := DB.QueryRow("select email, name from users where user_id=?",user_id).Scan(&user.Email, &user.Name)
	
	if err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(struct {User LoginUser} {user})
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, AUTH_COOKIE)

	session.Values[AUTH_COOKIE] = -1
	session.Save(r, w)
	
	fmt.Fprintln(w, "Logged out successfully!")
}

