package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)


type Video struct {
	VideoId string `db:"VideoId"`
	Title sql.NullString `db:"Title"`
	Description sql.NullString `db:"Description"`
	Year int `db:"Year"`
	ImgPath sql.NullString `db:"ImgPath"`
}

func GetVideos() []Video {
	db := GetDBContext()
	sql := "SELECT * FROM videos;"
	videos := []Video{}
	err := db.Select(&videos, sql)
	if err != nil {
		panic(err)
	}

	db.Close()

	return videos
	/*
	return []Video{
		{videoId: "123abc", title: "Jaws", description: "Harrowing Scary Movie"},
		{videoId: "456def", title: "Psycho", description: "Overrated ahhhh"},
		{videoId: "789ghi", title: "THe Good, the Bad and the Ugly", description: "Clint Eastwood is a terrible actor"},
	}
	*/
}


func GetVideoMap() map[string][]map[string]string {
	return map[string][]map[string]string{
		"videos": []map[string]string{
			{"videoId": "123abc", "title": "Jaws", "description": "Harrowing Scary Movie"},
			{"videoId": "456def", "title": "Psycho", "description": "Overrated ahhhh"},
			{"videoId": "789ghi", "title": "THe Good, the Bad and the Ugly", "description": "Clint Eastwood is a terrible actor"},
		},
	}
}
