package main

import (
	"log"
	"net/http"
	"westflix/controllers"
	"westflix/models"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()


	r.HandleFunc("/", controllers.ListVideos)

	r.HandleFunc("/register", func(resp http.ResponseWriter, req *http.Request) {
		controllers.PopulateViewBag(req)
		if req.Method == http.MethodPost {
			controllers.PostRegister(resp, req)
		} else if req.Method == http.MethodGet {
			controllers.GetRegister(resp, req)
		}
		models.CloseDB()
	})
	r.HandleFunc("/login", func(resp http.ResponseWriter, req *http.Request) {
		controllers.PopulateViewBag(req)
		if req.Method == http.MethodPost {
			controllers.PostLogin(resp, req)
		} else if req.Method == http.MethodGet {
			controllers.GetLogin(resp, req)
		}
		models.CloseDB()
	})
	r.HandleFunc("/logout", func(resp http.ResponseWriter, req *http.Request) {
		controllers.PopulateViewBag(req)
		if req.Method == http.MethodGet {
			controllers.PerformLogout(resp, req)
		}
		models.CloseDB()
	})

	r.HandleFunc("/watch/{videoId}/", controllers.ViewVideo)

	fs := http.FileServer(http.Dir("static/"))
	r.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":80", r))
}
