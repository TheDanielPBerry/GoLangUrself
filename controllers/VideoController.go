package controllers

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"westflix/models"
	"html/template"
)

type VideoList struct {
	videos []models.Video
}

func ListVideos(resp http.ResponseWriter, req *http.Request) {
	videos := models.GetVideos()

	tmpl, err := template.New("").ParseFiles("views/index.html", "views/base.html")
	if(err != nil) {
		panic(err)
	}
	fmt.Println(videos)
	tmpl.ExecuteTemplate(resp, "base", videos)
}

func ViewVideo(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	videoId := vars["videoId"]
	fmt.Fprintf(resp, "<video src=\"%s\"></video>", videoId)
}

