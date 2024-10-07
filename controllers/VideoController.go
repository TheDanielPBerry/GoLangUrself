package controllers

import (
	"net/http"
	"westflix/models"
	"strconv"
	"github.com/gorilla/mux"
	"encoding/json"
	"log"
	"html/template"
)


func ListVideos(resp http.ResponseWriter, req *http.Request) {
	tmpl := getTemplate("index")
	tmpl, err := tmpl.ParseFiles("views/row.html")
	if err != nil {
		log.Print(err)
	}

	mostPopular := models.GetMostPopularVideos()
	var popularRightJoin []models.Video
	
	session := models.GetSession(req)
	userId, ok := session.Values["UserId"].(int)
	if ok && userId > 0 {
		recentlyWatched := models.GetRecentlyWatchedVideos(userId)
		for _, popular := range *mostPopular {
			found := false
			for _, recent := range *recentlyWatched {
				if(popular.VideoId == recent.VideoId) {
					found = true
					break
				}
			}
			if !found {
				popularRightJoin = append(popularRightJoin, popular)
			}
		}
		viewBag["recentlyWatched"] = recentlyWatched
	}
	viewBag["mostPopular"] = popularRightJoin

	tmpl.ExecuteTemplate(resp, "base", viewBag)
}


func GetVideoPreview(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	vid := vars["videoId"]
	
	videoId, err := strconv.ParseInt(vid, 10, 64)
	if err != nil {
		http.Error(resp, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	video, ok := models.GetVideo(int(videoId))
	if !ok {
		http.Error(resp, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	tmpl, err := template.New("").ParseFiles("views/preview.html")
	if err != nil {
		log.Print(err)
		http.Error(resp, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	tmpl.ExecuteTemplate(resp, "preview", video)
}


func ViewVideo(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	vid := vars["videoId"]
	
	videoId, err := strconv.ParseInt(vid, 10, 64)
	if err != nil {
		http.Error(resp, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	video, ok := models.GetVideo(int(videoId))
	if !ok {
		http.Error(resp, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	viewBag["video"] = video


	session := models.GetSession(req)
	userId, ok := session.Values["UserId"].(int)
	if !ok || userId <= 0 {
		//Is not properly authenticated
		http.Redirect(resp, req, "/login", http.StatusSeeOther)
		return
	}

	watchEvent, ok := models.GetWatchEvent(userId, int(videoId))
	viewBag["watchEvent"] = watchEvent

	jsonData, err := json.Marshal(viewBag)
	if err != nil {
		http.Error(resp, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	viewBag["jsonData"] = string(jsonData)

	tmpl := getTemplate("watch")
	tmpl.ExecuteTemplate(resp, "base", viewBag)
}


func RecordWatchEvent(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	vid := vars["videoId"]
	videoId, err := strconv.ParseInt(vid, 10, 64)
	if err != nil || videoId < 0 || videoId > 101 {
		http.Error(resp, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	session := models.GetSession(req)
	userId, ok := session.Values["UserId"].(int)
	if !ok || userId <= 0 {
		//Is not properly authenticated
		json.NewEncoder(resp).Encode(map[string]interface{} {
			"error": "Invalid User",
			"errorCode": 403,
		})
		return
	}

	prog := vars["progress"]
	progress, err := strconv.ParseInt(prog, 10, 64)
	if err != nil {
		http.Error(resp, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	
	watchEvent, ok := models.GetWatchEvent(int(userId), int(videoId))
	if !ok || watchEvent == nil {
		watchEvent = new(models.WatchEvent)
		watchEvent.WatchEventId = 0
		watchEvent.UserId = int(userId)
		watchEvent.VideoId = int(videoId)
	}
	watchEvent.ProgressSeconds = int(progress)
	we, ok := models.UpdateWatchEvent(watchEvent)
	if !ok || we == nil {
		http.Error(resp, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	json.NewEncoder(resp).Encode(we)
}
