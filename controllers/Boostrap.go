package controllers


import (
	"net/http"
	"fmt"
	"log"
	"html/template"
	"westflix/models"
)

type ViewBag map[string]interface{}


func PopulateViewBag(req *http.Request) ViewBag {
	var viewBag ViewBag
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
	return viewBag
}

func (viewBag ViewBag) getTemplate(view string) *template.Template {
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
