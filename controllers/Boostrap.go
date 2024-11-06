package controllers


import (
	"net/http"
	"fmt"
	"log"
	"html/template"
	"westflix/models"
	"github.com/gorilla/sessions"
)

type ViewBag map[string]interface{}

type Response struct {
	w http.ResponseWriter
	viewBag ViewBag
	session *sessions.Session
	user *models.User
	Authenticated bool
}

func StartResponse(w http.ResponseWriter, req *http.Request) Response {
	var resp Response
	resp.w = w
	resp.viewBag = make(ViewBag)
	resp.session = models.GetSession(req)

	userId, ok := resp.session.Values["UserId"].(int)
	resp.viewBag["authenticated"] = false;
	resp.Authenticated = false
	if ok {
		resp.viewBag["authenticated"] = userId > 0
		resp.Authenticated = userId > 0
		user, ok := models.GetUserById(userId)
		if ok {
			resp.user = user
			frontFacingUser := models.User {
				UserId: userId,
				FullName: user.FullName,
				DateAdded: user.DateAdded,
				DateModified: user.DateModified,
			}
			resp.viewBag["user"] = frontFacingUser
		}
	}
	return resp
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
