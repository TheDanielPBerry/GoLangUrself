package main

import (
	"log"
	"net/http"
	"westflix/controllers"
	//"westflix/models"
	"github.com/gorilla/mux"
)
type Middleware func(http.HandlerFunc) http.HandlerFunc


type RouteFunc func(resp controllers.Response, req *http.Request)


func PrepRoute(f RouteFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := controllers.StartResponse(w, r)
		f(resp, r)
	}
}

func NoAuthAllowed(f RouteFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := controllers.StartResponse(w, r)
		if resp.Authenticated {
			http.Redirect(w, r, "/", http.StatusFound)
		}
		f(resp, r)
		return
	}
}

func RequireAuth(f RouteFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := controllers.StartResponse(w, r)
		if !resp.Authenticated {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		f(resp, r)
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

	r.HandleFunc("/v/{videoId}/preview", PrepRoute(controllers.GetVideoPreview)).Methods("GET")
	r.HandleFunc("/v/{videoId}/watch/{progress}", RequireAuth(controllers.RecordWatchEvent)).Methods("POST", "PUT")
	r.HandleFunc("/v/{videoId}/rate/{rating}", RequireAuth(controllers.RecordRating)).Methods("POST", "PUT")
	r.HandleFunc("/v/{videoId}", RequireAuth(controllers.ViewVideo)).Methods("GET")

	fs := http.FileServer(http.Dir("assets/"))
	r.Handle("/static/js/{$}", http.StripPrefix("/static/", fs)).Methods("GET")
	r.Handle("/static/css/{$}", http.StripPrefix("/static/", fs)).Methods("GET")
	r.Handle("/static/thumbnails/{$}", http.StripPrefix("/static/", fs)).Methods("GET")
	r.Handle("/static/videos/{$}", http.StripPrefix("/static/", fs)).Methods("GET")

	r.HandleFunc("/", PrepRoute(controllers.ListVideos)).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", r))
}
