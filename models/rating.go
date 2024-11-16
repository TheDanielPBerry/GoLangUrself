package models

import (
	"log"
)


type Rating struct {
	RatingId int `db:"RatingId"`
	UserId int `db:"UserId"`
	VideoId int `db:"VideoId"`
	Value int `db:"Value"`
	DateAdded string `db:"DateAdded"`
	DateModified string `db:"DateModified"`
}

func GetRating(userId int, videoId int) (*Rating, bool) {
	db := GetDBContext()

	sql := `SELECT RatingId, UserId, VideoId, 
	Value, DateAdded, DateModified
	FROM Rating
	WHERE UserId=? AND VideoId=?;`

	stmt, err := db.Preparex(sql)
	if err != nil {
		log.Fatal(err)
		return nil, false
	}

	var rating *Rating = new(Rating)
	err = stmt.Get(rating, userId, videoId)
	if err != nil {
		emptyRating := new(Rating)
		emptyRating.Value = 0
		emptyRating.UserId = userId
		emptyRating.VideoId = videoId
		return emptyRating, false
	}

	return rating, true
}

func UpdateRating(rating *Rating) (*Rating, bool) {
	var sql string
	if rating.RatingId <= 0 {
		sql = `INSERT INTO Rating
			(VideoId, UserId, Value, DateAdded, DateModified)
			VALUES
			(?,?,?,datetime(),datetime());`;
	} else {
		sql = `UPDATE Rating
			SET Value=?,
			DateModified=datetime()
			WHERE UserId=? 
			AND VideoId=?;`;
	}

	db := GetDBContext()
	stmt, err := db.Preparex(sql)
	if err != nil {
		return rating, false
	}

	if rating.RatingId <= 0 {
		result, err := stmt.Exec(rating.VideoId, rating.UserId, rating.Value)
		if err != nil {
			log.Print(err)
			return nil, false
		}
		insertId, err := result.LastInsertId()
		if err != nil {
			return rating, false
		}
		rating.RatingId = int(insertId)

	} else {
		_, err := stmt.Exec(rating.Value, rating.UserId, rating.VideoId)
		if err != nil {
			log.Print(err)
			return rating, false
		}
	}
	return rating, true
}

