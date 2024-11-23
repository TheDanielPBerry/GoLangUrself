package controllers

import (
	//"database/sql"
	"encoding/json"
	"html"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"westflix/models"

	"github.com/gorilla/mux"
)


func ListVideos(resp Response, req *http.Request) {
	tmpl := resp.viewBag.getTemplate("index")
	tmpl, err := tmpl.ParseFiles("views/row.html")
	if err != nil {
		log.Print(err)
	}

	mostPopular := *models.GetMostPopularVideos()
	var popularRightJoin []models.Video
	
	session := models.GetSession(req)
	userId, ok := session.Values["UserId"].(int)
	if ok && userId > 0 {
		recentlyWatched := models.GetRecentlyWatchedVideos(userId)
		for _, popular := range mostPopular {
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
		resp.viewBag["mostPopular"] = popularRightJoin[:10]

		if len(*recentlyWatched) > 0 {
			resp.viewBag["recentlyWatched"] = recentlyWatched

			recommendations := models.GetRecommendations(userId)
			if len(*recommendations) > 0 {
				resp.viewBag["forYou"] = recommendations
			}
		}

		queuedVideos := models.GetQueuedVideos(userId)
		if len(*queuedVideos) > 0 {
			resp.viewBag["backlog"] = queuedVideos
		}
	} else {
		resp.viewBag["mostPopular"] = mostPopular[0:10]
	}


	genres := models.GetRandomGenres()
	genreCollections := map[string]*[]models.Video{}
	for _, genre := range genres {
		genreCollections[genre.Description] = models.GetGenreVideos(genre.GenreId)
	}
	resp.viewBag["genres"] = genreCollections

	tmpl.ExecuteTemplate(resp.w, "base", resp.viewBag)
}


func GetVideoPreview(resp Response, req *http.Request) {
	vars := mux.Vars(req)
	vid := vars["videoId"]
	
	videoId, err := strconv.ParseInt(vid, 10, 64)
	if err != nil {
		http.Error(resp.w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	video, ok := models.GetVideo(int(videoId))
	if !ok {
		http.Error(resp.w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	tmpl, err := template.New("").ParseFiles("views/preview.html")
	if err != nil {
		log.Print(err)
		http.Error(resp.w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	video.Description.String = html.UnescapeString(video.Description.String)
	video.Queued = false

	session := models.GetSession(req)
	userId, ok := session.Values["UserId"].(int)
	if ok && userId > 0 {
		queueEvent, ok := models.GetWatchQueue(userId, int(videoId))
		if ok && queueEvent != nil {
			video.Queued = true
		}
	}

	tmpl.ExecuteTemplate(resp.w, "preview", video)
}


func ViewVideo(resp Response, req *http.Request) {
	vars := mux.Vars(req)
	vid := vars["videoId"]
	
	videoId, err := strconv.ParseInt(vid, 10, 64)
	if err != nil {
		http.Error(resp.w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	video, ok := models.GetVideo(int(videoId))
	if !ok {
		http.Error(resp.w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	resp.viewBag["video"] = video


	session := models.GetSession(req)
	userId, ok := session.Values["UserId"].(int)
	if !ok || userId <= 0 {
		//Is not properly authenticated
		http.Redirect(resp.w, req, "/login", http.StatusSeeOther)
		return
	}

	watchEvent, ok := models.GetWatchEvent(userId, int(videoId))
	resp.viewBag["watchEvent"] = watchEvent

	rating, ok := models.GetRating(userId, int(videoId));
	resp.viewBag["rating"] = rating;

	jsonData, err := json.Marshal(resp.viewBag)
	if err != nil {
		http.Error(resp.w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	resp.viewBag["jsonData"] = string(jsonData)

	tmpl := resp.viewBag.getTemplate("watch")
	tmpl, err = tmpl.ParseFiles("views/rating.html")
	if err != nil {
		log.Print(err)
	}
	tmpl.ExecuteTemplate(resp.w, "base", resp.viewBag)
}


func RecordWatchEvent(resp Response, req *http.Request) {
	vars := mux.Vars(req)

	vid := vars["videoId"]
	videoId, err := strconv.ParseInt(vid, 10, 64)
	if err != nil || videoId < 0 || videoId > 101 {
		http.Error(resp.w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	session := models.GetSession(req)
	userId, ok := session.Values["UserId"].(int)
	if !ok || userId <= 0 {
		//Is not properly authenticated
		json.NewEncoder(resp.w).Encode(map[string]interface{} {
			"error": "Invalid User",
			"errorCode": 403,
		})
		return
	}

	prog := vars["progress"]
	progress, err := strconv.ParseInt(prog, 10, 64)
	if err != nil {
		http.Error(resp.w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
		http.Error(resp.w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	json.NewEncoder(resp.w).Encode(we)
}

func SearchVideos(resp Response, req *http.Request) {
	vars := mux.Vars(req)

	query := vars["q"]
	videos := models.SearchVideos(query)
	if len(*videos) <= 0 {
		json.NewEncoder(resp.w).Encode([]map[string]interface{}{{
			"VideoId": "-1",
			"Title": map[string]interface{}{"String": "No Results"},
		}})
		return
	}
	json.NewEncoder(resp.w).Encode(videos)
}

func RecordRating(resp Response, req *http.Request) {
	vars := mux.Vars(req)

	vid := vars["videoId"]
	videoId, err := strconv.ParseInt(vid, 10, 64)
	if err != nil || videoId < 0 || videoId > 101 {
		http.Error(resp.w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	session := models.GetSession(req)
	userId, ok := session.Values["UserId"].(int)
	if !ok || userId <= 0 {
		//Is not properly authenticated
		json.NewEncoder(resp.w).Encode(map[string]interface{} {
			"error": "Invalid User",
			"errorCode": 403,
		})
		return
	}

	ratingInput := vars["rating"]
	ratingValue, err := strconv.ParseInt(ratingInput, 10, 64)
	if err != nil {
		http.Error(resp.w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	
	rating, ok := models.GetRating(int(userId), int(videoId))
	if !ok || rating == nil {
		rating = new(models.Rating)
		rating.Value = 0
		rating.UserId = int(userId)
		rating.VideoId = int(videoId)
	}
	rating.Value = int(ratingValue)
	we, ok := models.UpdateRating(rating)
	if !ok || we == nil {
		http.Error(resp.w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	json.NewEncoder(resp.w).Encode(we)
}


func QueueVideo(resp Response, req *http.Request) {
	vars := mux.Vars(req)

	vid := vars["videoId"]
	videoId, err := strconv.ParseInt(vid, 10, 64)
	if err != nil || videoId < 0 || videoId > 101 {
		http.Error(resp.w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	session := models.GetSession(req)
	userId, ok := session.Values["UserId"].(int)
	if !ok || userId <= 0 {
		//Is not properly authenticated
		json.NewEncoder(resp.w).Encode(map[string]interface{} {
			"error": "Invalid User",
			"errorCode": 403,
		})
		return
	}

	newQueueEvent := false
	queueEvent, ok := models.GetWatchQueue(int(userId), int(videoId))
	if !ok || queueEvent == nil {
		queueEvent = new(models.WatchQueue)
		queueEvent.WatchQueueId = 0
		queueEvent.UserId = int(userId)
		queueEvent.VideoId = int(videoId)
		newQueueEvent = true;

	}
	we, ok := models.UpdateWatchQueue(queueEvent)

	if !ok || we == nil {
		http.Error(resp.w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	if(newQueueEvent) {
		resp.w.WriteHeader(http.StatusCreated)
	} else {
		resp.w.WriteHeader(http.StatusNoContent)
	}
	json.NewEncoder(resp.w).Encode(we)
}
