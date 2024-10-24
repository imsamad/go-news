package main

import (	
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"regexp"
	"strings"
	"crypto/rand"
	"fmt"
	"net/http"
	"errors"
	"html/template"
	"os"
)

type TMiddleware func(http.HandlerFunc) http.HandlerFunc

func middlewares(f http.HandlerFunc, ms ...TMiddleware) http.HandlerFunc  {
	for _, m := range ms {
		f = m(f)
	}
	return f
}

func fetchAuthData(r *http.Request) AuthData {
	user, err := fetchSession(r)
	
	if  err != nil{
		return AuthData{
			IsLoggedIn:false,
			IsAdmin:false,
		}
	}
	
	return AuthData{
		IsLoggedIn:true,
		IsAdmin:user.Role == "ADMIN",
	}
}

func fetchSession (r *http.Request) (User, error)  {
	// if authMiddleware already called, on calling this fun reading from context
	// would avoid re-query user from db

	userPointer, ok := r.Context().Value(AUTH_COOKIE).(*User)
	
	if  ok && userPointer != nil {
		return *userPointer, nil
	}
	
	session, _ := store.Get(r, AUTH_COOKIE)
	user_id, ok := session.Values[AUTH_COOKIE].(int64)

	if user_id == -1 || !ok {
		return User{}, errors.New("not logged in")
	} 
	
	var user User;

	err := DB.QueryRow("select user_id,email,name,role from users where user_id=?",user_id).Scan(&user.UserId,&user.Email,&user.Name,&user.Role)

	if err != nil {
		return User{}, errors.New("not logged in")
	}

	return user, nil
}

func mustParseTemplate(filepath string) *template.Template {
    tmpl, err := template.ParseFiles(filepath)
    if err != nil {
		fmt.Println("reason being unable to load template file: ",err)
		os.Exit(1)
	}
    return tmpl
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func Slugify(s string) string {
	s = strings.ToLower(s)

	s = strings.ReplaceAll(s, " ", "-")

	re := regexp.MustCompile(`[^a-z0-9-]+`)
	s = re.ReplaceAllString(s, "")

	s = strings.Trim(s, "-")

	randomString, err := generateRandomString(5)
	if err != nil {
		return s
	}

	slug := fmt.Sprintf("%s-%s", s, randomString)

	return slug
}

func generateRandomString(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result strings.Builder

	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		result.WriteByte(letters[num.Int64()])
	}

	return result.String(), nil
}
