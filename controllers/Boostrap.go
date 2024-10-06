package controllers


import (
	"net/http"
	"fmt"
	"log"
	"html/template"
	"westflix/models"
)

type ViewBag map[string]interface{}

var viewBag ViewBag

func PopulateViewBag(req *http.Request) {
	viewBag = make(ViewBag)

	session := models.GetSession(req)
	userId, ok := session.Values["UserId"].(int)
	viewBag["authenticated"] = false;
	if ok {
		viewBag["authenticated"] = userId > 0
		user, ok := models.GetUserById(userId)
		if ok {
			viewBag["user"] = user
		}
	}
}

func getTemplate(view string) *template.Template {
	tmpl, err := template.New("").ParseFiles(fmt.Sprintf("views/%s.html", view), "views/base.html")
	if(err != nil) {
		panic(err)
	}
	
	headerRightViewPath := "views/header/login.html"
	if authenticated, ok := viewBag["authenticated"].(bool); ok && authenticated {
		headerRightViewPath = "views/header/account.html"
	}
	tmpl, err = tmpl.ParseFiles(headerRightViewPath)
	if err != nil {
		log.Panic(err)
	}

	return tmpl
}