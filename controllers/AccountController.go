package controllers

import (
	"net/http"
	"html/template"
	"westflix/validate"
	"westflix/models"
	//"fmt"
)

type PostBack struct {
	PostBack map[string]string
}

func GetLoginTemplate() *template.Template {
	tmpl, err := template.New("").ParseFiles("views/login.html", "views/base.html")
	if(err != nil) {
		panic(err)
	}
	return tmpl
}

func GetLogin(resp http.ResponseWriter, req *http.Request) {
	tmpl := GetLoginTemplate()
	var errList PostBack
	tmpl.ExecuteTemplate(resp, "base", errList)
}

func PostLogin(resp http.ResponseWriter, req *http.Request) {
	email := req.FormValue("email")
	if err := validate.Email(email); err != nil {
		tmpl := GetLoginTemplate()
		postBack := PostBack{map[string]string{"error": err.Error(), "email": email}}
		tmpl.ExecuteTemplate(resp, "base", postBack)
		return
	}

	password := req.FormValue("password")
	if err := validate.Password(password); err != nil {
		tmpl := GetLoginTemplate()
		postBack := PostBack{map[string]string{"error": err.Error(), "email": email}}
		tmpl.ExecuteTemplate(resp, "base", postBack)
		return
	}

	models.AttemptLogin(email, password)

}
