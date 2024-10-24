package main

import (
	"net/http"
)

var (
	adminPanelPageTmpl = mustParseTemplate("views/me.tmpl")
)

func getAdminPage(w http.ResponseWriter, r *http.Request) {
	user,_ := fetchSession(r)

	adminPanelPageTmpl.Execute(w, struct{IsAdminPage bool
		AuthData AuthData 
	   Posts []Post}{AuthData:fetchAuthData(r),Posts:getAllPostsExceptUtil(user.UserId),IsAdminPage:true})
}