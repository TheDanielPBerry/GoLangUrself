package models

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
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


func SearchVideos(query string) *[]Video {
	db := GetDBContext()
	sql := `SELECT VideoId, Title, Year FROM (
			SELECT *, 1 as priority
			FROM Video
			WHERE Title LIKE(?)
		UNION ALL
			SELECT *, 2 as priority
			FROM Video
			WHERE CONCAT(Title, Year) LIKE(?)
		UNION ALL
			SELECT *, 3 AS priority
			FROM Video
			WHERE Description LIKE(?)
		)
		GROUP BY VideoId
		ORDER BY priority, Title
		LIMIT 10
	`;
	stmt, err := db.Preparex(sql)
	if err != nil {
		log.Print(err)
		return nil
	}

	videos := new([]Video)

	query = strings.ReplaceAll(query, " ", "%");
	queryLead := query + "%"
	queryWildcard := "%" + query + "%"

	err = stmt.Select(videos, queryLead, queryWildcard, queryWildcard)

	if err != nil {
		log.Print(err)
		return nil
	}

	return videos
}


func GetMostPopularVideos() *[]Video {
	sql := `SELECT * FROM (
		SELECT awt.VideoId
		FROM AverageWatchTimes awt
		ORDER BY awt.TotalWatchTime DESC
		LIMIT 40
	)
	ORDER BY RANDOM()
	LIMIT 20`;

	db := GetDBContext()
	
	videos := new([]Video)
	err := db.Select(videos, sql)
	if err != nil {
		panic(err)
	}

	db.Close()

	return videos
}


func GetRecentlyWatchedVideos(userId int) *[]Video {
	sql := `SELECT we.VideoId
		FROM WatchEvent we
		WHERE we.UserId=?
		ORDER BY we.DateModified DESC
		LIMIT 10`;

	db := GetDBContext()
	stmt, err := db.Preparex(sql)
	if err != nil {
		log.Print(err)
		return nil
	}

	videos := new([]Video)
	err = stmt.Select(videos, userId)
	if err != nil {
		log.Print(err)
		return nil
	}

	return videos
}

func GetGenreVideos(genreId int) *[]Video {
	db := GetDBContext()
	sql := `SELECT vg.VideoId AS VideoId
		FROM VideoGenre vg
		WHERE vg.GenreId=?
		ORDER BY RANDOM()
		LIMIT 20;`
	stmt, err := db.Preparex(sql)
	if err != nil {
		log.Print(err)
		return nil
	}

	videos := new([]Video)
	err = stmt.Select(videos, genreId)
	if err != nil {
		log.Print(err)
		return nil
	}
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



