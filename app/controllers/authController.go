package controllers

import (
	"encoding/json"
	"fmt"
	"go-news/consts"
	"go-news/database"
	"go-news/lib"
	"go-news/session"
	"go-news/types"
	"go-news/views"
	"net/http"
)

type LoginPageData struct {
	Success string
	Fail string
	EmailMessage string
	PasswordMessage string
	AuthData types.AuthData
}

type SignupPageData struct {
	Success string
	Fail string
	EmailMessage string
	PasswordMessage string
	NameMessage string
	AuthData types.AuthData
}

 
// Load all templates globally and handle errors


func UserPanelPage(w http.ResponseWriter, r *http.Request) {

	user, _ := r.Context().Value(consts.AUTH_COOKIE).(*types.User)

	posts := lib.GetAllMyPostsUtil(user.UserId)

	views.UserPanelPageTmpl.Execute(w, struct{IsAdminPage bool
		 AuthData types.AuthData 
		Posts []types.Post}{AuthData:lib.FetchAuthData(r),Posts:posts,IsAdminPage:false})


}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	views.LoginPageTmpl.Execute(w, LoginPageData {
		Success: "",
		Fail :"",
		EmailMessage :"",
		PasswordMessage:"",
		AuthData:types.AuthData{},
	})
}

func SignupPage(w http.ResponseWriter, r *http.Request) {
	views.SignupPageTmpl.Execute(w, SignupPageData {
		Success: "",
		Fail :"",
		EmailMessage :"",
		PasswordMessage:"",
		NameMessage:"",
		AuthData:types.AuthData{},

	})
}

func Signup(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	name := r.FormValue("name")
	password := r.FormValue("password")
	
 
	if email == "" {
		views.SignupPageTmpl.Execute(w, SignupPageData {
			Success: "",
			Fail :"",
			EmailMessage : "Email is required!",
			PasswordMessage: "",
			NameMessage:"",			AuthData:types.AuthData{},

		})
		return
	} 

	if password == "" {
		views.SignupPageTmpl.Execute(w, SignupPageData {
			Success: "",
			Fail :"",
			EmailMessage : "",
			PasswordMessage:"Password is required",
			NameMessage:"",			AuthData:types.AuthData{},

		})
		return
		} 

		if name == "" {
			views.SignupPageTmpl.Execute(w, SignupPageData {
				Success: "",
				Fail :"",
				EmailMessage : "",
				PasswordMessage:"",
				NameMessage:"Name is required",			AuthData:types.AuthData{},

			})
			return
			} 
		

	alreadyExistedQuery := `
		select email from users where email=?
	`

	notExist := database.GetDB().QueryRow(alreadyExistedQuery, email).Scan(&email)
	if notExist == nil {
		views.SignupPageTmpl.Execute(w, SignupPageData {
			Success: "",
			Fail :"",
			EmailMessage : "Email already exist!",
			PasswordMessage: "",
			NameMessage:"",			AuthData:types.AuthData{},

		})
		return
	}	

	userQuery := `
		insert into users (name, email,password) values (?,?,?);
	`
	
    hashPwd, _ := lib.HashPassword(password)		
	newUser, err := database.GetDB().Exec(userQuery, name, email, hashPwd)


	if err != nil {
		views.SignupPageTmpl.Execute(w, SignupPageData {
			Success: "",
			Fail :"Plz try again",
			EmailMessage : "",
			PasswordMessage: "",
			NameMessage:"",			AuthData:types.AuthData{},

		})
		return
	}	
	
	id, err := newUser.LastInsertId()
	
	if err != nil {

		views.SignupPageTmpl.Execute(w, SignupPageData {
			Success: "",
			Fail :"Plz try again",
			EmailMessage : "",
			PasswordMessage: "",
			NameMessage:"",			AuthData:types.AuthData{},

		})
		return
	}
	
	session, _ := session.Store.Get(r, consts.AUTH_COOKIE)

	session.Values[consts.AUTH_COOKIE] = id
	session.Save(r, w)
	
	http.Redirect(w,r,"/me",http.StatusFound)
}

func HandleLogin(w http.ResponseWriter, r * http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" {
		views.LoginPageTmpl.Execute(w, LoginPageData {
			Success: "",
			Fail :"",
			EmailMessage : "Email is required!",
			PasswordMessage: "",
			AuthData:types.AuthData{},

		})
		return
	} 

	if password == "" {
		views.LoginPageTmpl.Execute(w, LoginPageData {
			Success: "",
			Fail :"",
			EmailMessage : "",
			PasswordMessage:"Password is required",
			AuthData:types.AuthData{},

		})
		return
		} 
	
	var dbPwd string;
	var user_id int64;
	var name string;
	var role types.ROLE;

	err := database.GetDB().QueryRow("select user_id, email, name, password,role from users where email=?",email).Scan(&user_id, &email, &name, &dbPwd,&role)

	if err != nil {
		views.LoginPageTmpl.Execute(w, LoginPageData {
			Success: "",
			Fail :"Email does not exist",
			EmailMessage : "",
			PasswordMessage: "",
			AuthData:types.AuthData{},

		})
		return
	}

	
	if !lib.CheckPasswordHash(password, dbPwd) {
		views.LoginPageTmpl.Execute(w, LoginPageData {
			Success: "",
			Fail :"",
			EmailMessage : "",
			PasswordMessage: "Password is incorrect!",
			AuthData:types.AuthData{},

		})
		return
	}

	_ = types.User {
		UserId:user_id,
		Name:name,
		Email:email,
		Role:role,		
	}

	session , _ := session.Store.Get(r, consts.AUTH_COOKIE)
 
	session.Values[consts.AUTH_COOKIE] = user_id
	session.Save(r,w)
	// json.NewEncoder(w).Encode(user)
	http.Redirect(w,r,"/me",http.StatusFound)

}


func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := session.Store.Get(r, consts.AUTH_COOKIE)

	session.Values[consts.AUTH_COOKIE] = -1
	session.Save(r, w)
	http.Redirect(w,r,"/auth/login",http.StatusFound)
	fmt.Fprintln(w, "Logged out successfully!")
}

func Me(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(consts.AUTH_COOKIE).(*types.User)

	if  !ok || user == nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(user)
}
