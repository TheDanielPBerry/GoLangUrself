package main

import (
	"log"
	"net/http"
	"westflix/controllers"
	//"westflix/models"
	"github.com/gorilla/mux"
)
type Middleware func(http.HandlerFunc) http.HandlerFunc



func PrepRoute(f http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		controllers.PopulateViewBag(req)
		f(resp, req)
	}
}

func NoAuthAllowed(f http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		viewBag := controllers.PopulateViewBag(req)
		if auth, ok := viewBag["authenticated"].(bool); ok && auth {
			http.Redirect(resp, req, "/", http.StatusFound)
		}
		f(resp, req)
		return
	}
}

func RequireAuth(f http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		viewBag := controllers.PopulateViewBag(req)
		if auth, ok := viewBag["authenticated"].(bool); !ok || !auth {
			http.Redirect(resp, req, "/login", http.StatusFound)
			return
		}
		f(resp, req)
	}
}

func main() {
	r := mux.NewRouter()


	r.HandleFunc("/register", NoAuthAllowed(controllers.GetRegister)).Methods("GET")
	r.HandleFunc("/register", NoAuthAllowed(controllers.PostRegister)).Methods("POST")

	r.HandleFunc("/login", NoAuthAllowed(controllers.GetLogin)).Methods("GET")
	r.HandleFunc("/login", NoAuthAllowed(controllers.PostLogin)).Methods("POST")

	r.HandleFunc("/login", PrepRoute(controllers.GetLogin)).Methods("GET")
	r.HandleFunc("/login", PrepRoute(controllers.PostLogin)).Methods("POST")
	r.HandleFunc("/logout", RequireAuth(controllers.PerformLogout)).Methods("GET")

	r.HandleFunc("/search/{q}", PrepRoute(controllers.SearchVideos)).Methods("GET")

	r.HandleFunc("/v/{videoId}/preview", controllers.GetVideoPreview).Methods("GET")
	r.HandleFunc("/v/{videoId}/watch/{progress}", RequireAuth(controllers.RecordWatchEvent)).Methods("POST", "PUT")
	r.HandleFunc("/v/{videoId}", RequireAuth(controllers.ViewVideo)).Methods("GET")

	fs := http.FileServer(http.Dir("assets/"))
	r.Handle("/static/js/{$}", http.StripPrefix("/static/", fs)).Methods("GET")
	r.Handle("/static/css/{$}", http.StripPrefix("/static/", fs)).Methods("GET")
	r.Handle("/static/thumbnails/{$}", http.StripPrefix("/static/", fs)).Methods("GET")
	r.Handle("/static/videos/{$}", http.StripPrefix("/static/", fs)).Methods("GET")

	r.HandleFunc("/", PrepRoute(controllers.ListVideos)).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", r))
}
