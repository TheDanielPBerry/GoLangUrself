package controllers

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"westflix/models"
)


func ListVideos(resp http.ResponseWriter, req *http.Request) {
	PopulateViewBag(req)

	tmpl := getTemplate("index")

	viewBag["videos"] = models.GetVideos()

	tmpl.ExecuteTemplate(resp, "base", viewBag)

}


func ViewVideo(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	videoId := vars["videoId"]
	fmt.Fprintf(resp, "<video src=\"%s\"></video>", videoId)
}

