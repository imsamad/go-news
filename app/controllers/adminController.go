package controllers

import (
	"go-news/lib"
	"go-news/types"
	"go-news/views"
	"net/http"
)
 

func GetAdminPage(w http.ResponseWriter, r *http.Request) {
	user,_ := lib.FetchSession(r)

	views.AdminPanelPageTmpl.Execute(w, struct{IsAdminPage bool
		AuthData types.AuthData 
	   Posts []types.Post}{AuthData:lib.FetchAuthData(r),Posts:lib.GetAllPostsExceptUtil(user.UserId),IsAdminPage:true})
}