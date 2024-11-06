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



func performLogin(id int, resp Response, req *http.Request) {
	session := models.GetSession(req)

	session.Values["UserId"] = id
	session.Save(req, resp.w)
}

func failLogin(resp Response) {
	tmpl := resp.viewBag.getTemplate("login")
	tmpl.ExecuteTemplate(resp.w, "base", resp.viewBag)
}

func failRegister(resp Response, err error, email string, fullName string) {
	tmpl := resp.viewBag.getTemplate("register")

	resp.viewBag["error"] = err.Error()
	resp.viewBag["email"] = email
	resp.viewBag["fullname"] = fullName

	tmpl.ExecuteTemplate(resp.w, "base", resp.viewBag)
}



/*********************
* Controller Actions *
*********************/

func GetLogin(resp Response, req *http.Request) {
	tmpl := resp.viewBag.getTemplate("login")
	tmpl.ExecuteTemplate(resp.w, "base", resp.viewBag)
}

func PostLogin(resp Response, req *http.Request) {
	email := req.FormValue("email")
	resp.viewBag["email"] = email
	if err := validate.Email(email); err != nil {
		resp.viewBag["error"] = "Invalid Credentials";
		failLogin(resp)
		return
	}

	password := req.FormValue("password")
	if err := validate.Password(password); err != nil {
		resp.viewBag["error"] = "Invalid Credentials";
		failLogin(resp)
		return
	}

	user, err := models.AuthenticateLogin(email, password)
	if err != nil {
		resp.viewBag["error"] = err.Error()
		failLogin(resp)
		return
	}

	performLogin(user.UserId, resp, req)

	http.Redirect(resp.w, req, "/", http.StatusSeeOther)
}




func GetRegister(resp Response, req *http.Request) {
	tmpl := resp.viewBag.getTemplate("register")
	tmpl.ExecuteTemplate(resp.w, "base", nil)
}


func PostRegister(resp Response, req *http.Request) {
	email := req.FormValue("email")
	fullName := req.FormValue("fullname")

	if err := validate.Email(email); err != nil {
		failRegister(resp, err, email, fullName)
		return
	}

	password := req.FormValue("password")
	if err := validate.Password(password); err != nil {
		failRegister(resp, err, email, fullName)
		return
	}

	if err:= validate.FullName(fullName); err != nil {
		failRegister(resp, err, email, fullName)
	}
	user := models.User{FullName: fullName, Email: email, Password: password}
	userId, err := models.CreateUser(user)
	if err != nil {
		failRegister(resp, err, email, fullName)
		return
	}
	performLogin(userId, resp, req)
	http.Redirect(resp.w, req, "/", http.StatusSeeOther)

	models.CloseDB()
}


func PerformLogout(resp Response, req *http.Request) {
	if auth, ok := resp.viewBag["authenticated"].(bool); auth && ok {
		session := models.GetSession(req)
		delete(session.Values, "UserId")
		session.Save(req, resp.w)
	}
	http.Redirect(resp.w, req, "/", http.StatusSeeOther)
}


