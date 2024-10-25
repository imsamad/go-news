package middlewares

import (
	"context"
	"go-news/consts"
	"go-news/lib"
	"net/http"
)


func AuthMiddleware() TMiddleware {
	return func (f http.HandlerFunc) http.HandlerFunc {
		return func (w http.ResponseWriter, r *http.Request) {
			user, err := lib.FetchSession(r)
			if  err != nil{
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			
			ctx := context.WithValue(r.Context(), consts.AUTH_COOKIE, &user)
		 
			f(w, r.WithContext(ctx))
		}	
	}
}


func AdminMiddleware() TMiddleware {
	return func (f http.HandlerFunc) http.HandlerFunc {
		return func (w http.ResponseWriter, r *http.Request) {
			user, err := lib.FetchSession(r)
			if  err != nil{
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			if user.Role != "ADMIN" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			
			ctx := context.WithValue(r.Context(), consts.AUTH_COOKIE, &user)
		 
			f(w, r.WithContext(ctx))
		}	
	}
}



func GuestMiddleware() TMiddleware {
	return func (f http.HandlerFunc) http.HandlerFunc {
		return func (w http.ResponseWriter, r *http.Request) {
			_, err := lib.FetchSession(r)
			
			if  err != nil{
				f(w,r)
				return
			}
			
			http.Redirect(w, r, "/me",http.StatusFound)
		}	
	}
}

