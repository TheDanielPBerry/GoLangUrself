package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)


type Video struct {
	VideoId string `db:"VideoId"`
	Title sql.NullString `db:"Title"`
	Description sql.NullString `db:"Description"`
	Year int `db:"Year"`
	TotalRuntime int `db:"TotalRuntime"`
	RuntimeDisplay sql.NullString `db:"RuntimeDisplay"`
	MPARating sql.NullString `db:"MPARating"`
	TotalRuntimeSeconds int `db:"TotalRuntimeSeconds"`
}


func GetVideos() []Video {
	db := GetDBContext()
	sql := `SELECT v.VideoId, v.Title, v.Description, v.Year,
		v.TotalRuntime, v.RuntimeDisplay, v.MPARating,
		(v.TotalRuntime * 60) AS TotalRuntimeSeconds
		FROM Video v
		ORDER BY RANDOM()
		LIMIT 10;`

	videos := []Video{}
	err := db.Select(&videos, sql)
	if err != nil {
		panic(err)
	}

	db.Close()

	return videos
}


func GetUnfinishedVideos(userId int) *[]Video {
	db := GetDBContext()
	sql := `SELECT v.VideoId
		FROM WatchEvent we
		INNER JOIN Video v ON v.VideoId=we.VideoID
		WHERE (v.TotalRuntime*60)>we.ProgressSeconds AND we.UserId=?
		ORDER BY we.DateModified DESC
		LIMIT 20;`
	stmt, err := db.Preparex(sql)
	if err != nil {
		log.Print(err)
		return nil
	}

	videos := new([]Video)
	err = stmt.Get(videos, userId)
	if err != nil {
		log.Print(err)
		return nil
	}

	return videos
}


func GetUserFavorites(userId int) *[]Video {
	db := GetDBContext()
	sql := `SELECT v.VideoId
		FROM Rating r
		INNER JOIN Video v ON v.VideoId=r.VideoID
		WHERE r.value>0 AND r.UserId=?
		ORDER BY r.DateAdded DESC
		LIMIT 20;`
	stmt, err := db.Preparex(sql)
	if err != nil {
		log.Print(err)
		return nil
	}

	videos := new([]Video)
	err = stmt.Get(videos, userId)
	if err != nil {
		log.Print(err)
		return nil
	}

	return videos
}


func GetVideo(id int) (*Video, bool) {
	db := GetDBContext()
	sql := `SELECT v.VideoId, v.Title, v.Description, v.Year,
		v.TotalRuntime, v.RuntimeDisplay, v.MPARating,
		(v.TotalRuntime * 60) AS TotalRuntimeSeconds
		FROM Video v
		WHERE v.VideoId=?`

	stmt, err := db.Preparex(sql)
	if err != nil {
		log.Print(err)
		return nil, false
	}

	var video *Video = new(Video)
	err = stmt.Get(video, id)
	if err != nil {
		log.Print(err)
		return nil, false
	}

	return video, true;
}


