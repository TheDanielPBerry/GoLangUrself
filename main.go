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

	r.HandleFunc("/v/{videoId}/watch/{progress}", controllers.RecordWatchEvent)
	r.HandleFunc("/v/{videoId}/preview", controllers.GetVideoPreview)
	r.HandleFunc("/v/{videoId}", controllers.ViewVideo)

	fs := http.FileServer(http.Dir("assets/"))
	r.Handle("/static/js/{$}", http.StripPrefix("/static/", fs))
	r.Handle("/static/css/{$}", http.StripPrefix("/static/", fs))
	r.Handle("/static/thumbnails/{$}", http.StripPrefix("/static/", fs))
	r.Handle("/static/videos/{$}", http.StripPrefix("/static/", fs))

	r.HandleFunc("/", controllers.ListVideos)

	log.Fatal(http.ListenAndServe(":80", r))
}
