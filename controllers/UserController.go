package controllers

import (
	//"errors"
	//"html/template"
	//"log"
	"net/http"
	"westflix/models"
	"westflix/validate"
	//"fmt"
)



func performLogin(id int, resp http.ResponseWriter, req *http.Request) {
	session := models.GetSession(req)

	session.Values["UserId"] = id
	session.Save(req, resp)
}

func failLogin(resp http.ResponseWriter, viewBag ViewBag) {
	tmpl := viewBag.getTemplate("login")
	tmpl.ExecuteTemplate(resp, "base", viewBag)
}

func failRegister(resp http.ResponseWriter, err error, email string, fullName string, viewBag ViewBag) {
	tmpl := viewBag.getTemplate("register")

	viewBag["error"] = err.Error()
	viewBag["email"] = email
	viewBag["fullname"] = fullName

	tmpl.ExecuteTemplate(resp, "base", viewBag)
}



/*********************
* Controller Actions *
*********************/

func GetLogin(resp http.ResponseWriter, req *http.Request, viewBag ViewBag) {
	tmpl := viewBag.getTemplate("login")
	tmpl.ExecuteTemplate(resp, "base", viewBag)
}

func PostLogin(resp http.ResponseWriter, req *http.Request, viewBag ViewBag) {
	email := req.FormValue("email")
	viewBag["email"] = email
	if err := validate.Email(email); err != nil {
		viewBag["error"] = "Invalid Credentials";
		failLogin(resp, viewBag)
		return
	}

	password := req.FormValue("password")
	if err := validate.Password(password); err != nil {
		viewBag["error"] = "Invalid Credentials";
		failLogin(resp, viewBag)
		return
	}

	user, err := models.AuthenticateLogin(email, password)
	if err != nil {
		viewBag["error"] = err.Error()
		failLogin(resp, viewBag)
		return
	}

	performLogin(user.UserId, resp, req)

	http.Redirect(resp, req, "/", http.StatusSeeOther)
}




func GetRegister(resp http.ResponseWriter, req *http.Request, viewBag ViewBag) {
	tmpl := viewBag.getTemplate("register")
	tmpl.ExecuteTemplate(resp, "base", nil)
}


func PostRegister(resp http.ResponseWriter, req *http.Request, viewBag ViewBag) {
	email := req.FormValue("email")
	fullName := req.FormValue("fullname")

	if err := validate.Email(email); err != nil {
		failRegister(resp, err, email, fullName, viewBag)
		return
	}

	password := req.FormValue("password")
	if err := validate.Password(password); err != nil {
		failRegister(resp, err, email, fullName, viewBag)
		return
	}

	if err:= validate.FullName(fullName); err != nil {
		failRegister(resp, err, email, fullName, viewBag)
	}
	user := models.User{FullName: fullName, Email: email, Password: password}
	userId, err := models.CreateUser(user)
	if err != nil {
		failRegister(resp, err, email, fullName, viewBag)
		return
	}
	performLogin(userId, resp, req)
	http.Redirect(resp, req, "/", http.StatusSeeOther)

	models.CloseDB()
}


func PerformLogout(resp http.ResponseWriter, req *http.Request, viewBag ViewBag) {
	if auth, ok := viewBag["authenticated"].(bool); auth && ok {
		session := models.GetSession(req)
		delete(session.Values, "UserId")
		session.Save(req, resp)
	}
	http.Redirect(resp, req, "/", http.StatusSeeOther)
}


