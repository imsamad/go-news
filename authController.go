package main

import (
	"net/http"
	"fmt"
	"encoding/json"
)

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
	 
	json.NewEncoder(w).Encode(User {Name:name, Email:email, UserId:id, Role: "USER"})
}

func handleLogin(w http.ResponseWriter, r * http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	
	if email == "" || password == "" {
		fmt.Fprintln(w, http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return
	} 
	
	var dbPwd string;
	var user_id int64;
	var name string;
	var role ROLE;

	err := DB.QueryRow("select user_id, email, name, password,role from users where email=?",email).Scan(&user_id, &email, &name, &dbPwd,&role)

	if err != nil {
		fmt.Println(err)
		return
	}

	
	if !CheckPasswordHash(password, dbPwd) {
		fmt.Println(err)
		return
	}

	user := User {
		UserId:user_id,
		Name:name,
		Email:email,
		Role:role,		
	}

	session , err := store.Get(r, AUTH_COOKIE)
 
	session.Values[AUTH_COOKIE] = user_id
	session.Save(r,w)
	json.NewEncoder(w).Encode(user)
}

func me(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(AUTH_COOKIE).(*User)

	if  !ok || user == nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, AUTH_COOKIE)

	session.Values[AUTH_COOKIE] = -1
	session.Save(r, w)
	
	fmt.Fprintln(w, "Logged out successfully!")
}

