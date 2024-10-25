package views

import (
	"fmt"
	"html/template"
	"os"
)

var (
    LoginPageTmpl    = MustParseTemplate("templates/login.tmpl")
    SignupPageTmpl   = MustParseTemplate("templates/signup.tmpl")
    UserPanelPageTmpl = MustParseTemplate("templates/me.tmpl")
	HomePageTmpl = MustParseTemplate("templates/index.tmpl")
	AdminPanelPageTmpl = MustParseTemplate("templates/me.tmpl")
    PostFormTmpl    = MustParseTemplate("templates/post_form.tmpl")
)

func MustParseTemplate(filepath string) *template.Template {
    tmpl, err := template.ParseFiles(filepath)
    if err != nil {
		fmt.Println("reason being unable to load template file: ",err)
		os.Exit(1)
	}
    return tmpl
}
