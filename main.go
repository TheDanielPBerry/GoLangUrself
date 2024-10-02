package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"westflix/controllers"
)

func main() {
	r := mux.NewRouter()


	r.HandleFunc("/", controllers.ListVideos)
	r.HandleFunc("/login", func(resp http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			controllers.PostLogin(resp, req)
		} else if req.Method == http.MethodGet {
			controllers.GetLogin(resp, req)
		}
	})

	r.HandleFunc("/watch/{videoId}/", controllers.ViewVideo)

	fs := http.FileServer(http.Dir("static/"))
	r.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":80", r))
}
