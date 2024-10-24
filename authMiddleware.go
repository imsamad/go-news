package main

import (
	"net/http"
	"context"
)



func authMiddleware() TMiddleware {
	return func (f http.HandlerFunc) http.HandlerFunc {
		return func (w http.ResponseWriter, r *http.Request) {
			user, err := fetchSession(r)
			if  err != nil{
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			
			ctx := context.WithValue(r.Context(), AUTH_COOKIE, &user)
		 
			f(w, r.WithContext(ctx))
		}	
	}
}


func adminMiddleware() TMiddleware {
	return func (f http.HandlerFunc) http.HandlerFunc {
		return func (w http.ResponseWriter, r *http.Request) {
			user, err := fetchSession(r)
			if  err != nil{
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			if user.Role != "ADMIN" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			
			ctx := context.WithValue(r.Context(), AUTH_COOKIE, &user)
		 
			f(w, r.WithContext(ctx))
		}	
	}
}

func isAdmin(r *http.Request) bool {
	user,ok := r.Context().Value(AUTH_COOKIE).(*User)
	if user == nil || !ok || user.Role != "ADMIN" {
		return false
	}	
	return true
}

func guestMiddleware() TMiddleware {
	return func (f http.HandlerFunc) http.HandlerFunc {
		return func (w http.ResponseWriter, r *http.Request) {
			_, err := fetchSession(r)
			
			if  err != nil{
				f(w,r)
				return
			}
			
			http.Redirect(w, r, "/me",http.StatusFound)
		}	
	}
}

