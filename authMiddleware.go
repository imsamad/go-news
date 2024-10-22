package main

import (
	"net/http"
	"context"
)

type TMiddleware func(http.HandlerFunc) http.HandlerFunc


func authMiddleware() TMiddleware {
	return func (f http.HandlerFunc) http.HandlerFunc {
		return func (w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, AUTH_COOKIE)
			user_id, ok := session.Values[AUTH_COOKIE].(int64)
	
			if  !ok || user_id == -1 {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			var user User;

			err:= DB.QueryRow("select user_id,email,name,role from users where user_id=?",user_id).Scan(&user.UserId,&user.Email,&user.Name,&user.Role)

			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), AUTH_COOKIE, &user)

		 
			f(w, r.WithContext(ctx))
		}	
	}
}

func middlewares(f http.HandlerFunc, ms ...TMiddleware) http.HandlerFunc  {
	for _, m := range ms {
		f = m(f)
	}
	return f
}
