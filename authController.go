package main

import (
	"net/http"
	"fmt"
	"encoding/json"
)

type LoginPageData struct {
	Success string
	Fail string
	EmailMessage string
	PasswordMessage string
	AuthData AuthData
}

type SignupPageData struct {
	Success string
	Fail string
	EmailMessage string
	PasswordMessage string
	NameMessage string
	AuthData AuthData
}

 
// Load all templates globally and handle errors
var (
    loginPageTmpl    = mustParseTemplate("views/login.tmpl")
    signupPageTmpl   = mustParseTemplate("views/signup.tmpl")
    userPanelPageTmpl = mustParseTemplate("views/me.tmpl")
)

func userPanelPage(w http.ResponseWriter, r *http.Request) {

	user, _ := r.Context().Value(AUTH_COOKIE).(*User)

	posts := getAllMyPostsUtil(user.UserId)

	userPanelPageTmpl.Execute(w, struct{IsAdminPage bool
		 AuthData AuthData 
		Posts []Post}{AuthData:fetchAuthData(r),Posts:posts,IsAdminPage:false})


}

func loginPage(w http.ResponseWriter, r *http.Request) {
	loginPageTmpl.Execute(w, LoginPageData {
		Success: "",
		Fail :"",
		EmailMessage :"",
		PasswordMessage:"",
		AuthData:AuthData{},
	})
}

func signupPage(w http.ResponseWriter, r *http.Request) {
	signupPageTmpl.Execute(w, SignupPageData {
		Success: "",
		Fail :"",
		EmailMessage :"",
		PasswordMessage:"",
		NameMessage:"",
		AuthData:AuthData{},

	})
}

func signup(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	name := r.FormValue("name")
	password := r.FormValue("password")
	
 
	if email == "" {
		signupPageTmpl.Execute(w, SignupPageData {
			Success: "",
			Fail :"",
			EmailMessage : "Email is required!",
			PasswordMessage: "",
			NameMessage:"",			AuthData:AuthData{},

		})
		return
	} 

	if password == "" {
		signupPageTmpl.Execute(w, SignupPageData {
			Success: "",
			Fail :"",
			EmailMessage : "",
			PasswordMessage:"Password is required",
			NameMessage:"",			AuthData:AuthData{},

		})
		return
		} 

		if name == "" {
			signupPageTmpl.Execute(w, SignupPageData {
				Success: "",
				Fail :"",
				EmailMessage : "",
				PasswordMessage:"",
				NameMessage:"Name is required",			AuthData:AuthData{},

			})
			return
			} 
		

	alreadyExistedQuery := `
		select email from users where email=?
	`

	notExist := DB.QueryRow(alreadyExistedQuery, email).Scan(&email)
	if notExist == nil {
		signupPageTmpl.Execute(w, SignupPageData {
			Success: "",
			Fail :"",
			EmailMessage : "Email already exist!",
			PasswordMessage: "",
			NameMessage:"",			AuthData:AuthData{},

		})
		return
	}	

	userQuery := `
		insert into users (name, email,password) values (?,?,?);
	`
	
    hashPwd, _ := HashPassword(password)		
	newUser, err := DB.Exec(userQuery, name, email, hashPwd)


	if err != nil {
		signupPageTmpl.Execute(w, SignupPageData {
			Success: "",
			Fail :"Plz try again",
			EmailMessage : "",
			PasswordMessage: "",
			NameMessage:"",			AuthData:AuthData{},

		})
		return
	}	
	
	id, err := newUser.LastInsertId()
	
	if err != nil {

		signupPageTmpl.Execute(w, SignupPageData {
			Success: "",
			Fail :"Plz try again",
			EmailMessage : "",
			PasswordMessage: "",
			NameMessage:"",			AuthData:AuthData{},

		})
		return
	}
	
	session, _ := store.Get(r, AUTH_COOKIE)

	session.Values[AUTH_COOKIE] = id
	session.Save(r, w)
	
	http.Redirect(w,r,"/me",http.StatusFound)
}

func handleLogin(w http.ResponseWriter, r * http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" {
		loginPageTmpl.Execute(w, LoginPageData {
			Success: "",
			Fail :"",
			EmailMessage : "Email is required!",
			PasswordMessage: "",
			AuthData:AuthData{},

		})
		return
	} 

	if password == "" {
		loginPageTmpl.Execute(w, LoginPageData {
			Success: "",
			Fail :"",
			EmailMessage : "",
			PasswordMessage:"Password is required",
			AuthData:AuthData{},

		})
		return
		} 
	
	var dbPwd string;
	var user_id int64;
	var name string;
	var role ROLE;

	err := DB.QueryRow("select user_id, email, name, password,role from users where email=?",email).Scan(&user_id, &email, &name, &dbPwd,&role)

	if err != nil {
		loginPageTmpl.Execute(w, LoginPageData {
			Success: "",
			Fail :"Email does not exist",
			EmailMessage : "",
			PasswordMessage: "",
			AuthData:AuthData{},

		})
		return
	}

	
	if !CheckPasswordHash(password, dbPwd) {
		loginPageTmpl.Execute(w, LoginPageData {
			Success: "",
			Fail :"",
			EmailMessage : "",
			PasswordMessage: "Password is incorrect!",
			AuthData:AuthData{},

		})
		return
	}

	_ = User {
		UserId:user_id,
		Name:name,
		Email:email,
		Role:role,		
	}

	session , err := store.Get(r, AUTH_COOKIE)
 
	session.Values[AUTH_COOKIE] = user_id
	session.Save(r,w)
	// json.NewEncoder(w).Encode(user)
	http.Redirect(w,r,"/me",http.StatusFound)

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
	http.Redirect(w,r,"/auth/login",http.StatusFound)
	fmt.Fprintln(w, "Logged out successfully!")
}

