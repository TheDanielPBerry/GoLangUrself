package models

import (
	"log"
)


type WatchEvent struct {
	WatchEventId int `db:"WatchEventId"`
	UserId int `db:"UserId"`
	VideoId int `db:"VideoId"`
	ProgressSeconds int `db:"ProgressSeconds"`
	DateAdded string `db:"DateAdded"`
	DateModified string `db:"DateModified"`
}

func GetWatchEvent(userId int, videoId int) (*WatchEvent, bool) {
	db := GetDBContext()

	sql := `SELECT we.WatchEventId, we.VideoId, we.UserId,
		we.ProgressSeconds, we.DateAdded, we.DateModified
		FROM WatchEvent we
		WHERE we.UserId=? AND we.VideoId=?;`

	stmt, err := db.Preparex(sql)
	if err != nil {
		log.Fatal(err)
		return nil, false
	}

	var watchEvent *WatchEvent = new(WatchEvent)
	err = stmt.Get(watchEvent, userId, videoId)
	if err != nil {
		log.Print(err)
		return nil, false
	}

	return watchEvent, true
}

func UpdateWatchEvent(watchEvent *WatchEvent) (*WatchEvent, bool) {
	var sql string
	if watchEvent.WatchEventId <= 0 {
		sql = `INSERT INTO WatchEvent
			(VideoId, UserId, ProgressSeconds, DateAdded, DateModified)
			VALUES
			(?,?,?,datetime(),datetime());`;
	} else {
		sql = `UPDATE WatchEvent
			SET ProgressSeconds=?,
			DateModified=datetime()
			WHERE UserId=? 
			AND VideoId=?;`;
	}

	db := GetDBContext()
	stmt, err := db.Preparex(sql)
	if err != nil {
		return watchEvent, false
	}

	if watchEvent.WatchEventId <= 0 {
		result, err := stmt.Exec(watchEvent.VideoId, watchEvent.UserId, watchEvent.ProgressSeconds)
		if err != nil {
			log.Print(err)
			return nil, false
		}
		insertId, err := result.LastInsertId()
		if err != nil {
			return watchEvent, false
		}
		watchEvent.WatchEventId = int(insertId)

	} else {
		_, err := stmt.Exec(watchEvent.ProgressSeconds, watchEvent.UserId, watchEvent.VideoId)
		if err != nil {
			log.Print(err)
			return watchEvent, false
		}
	}
	return watchEvent, true
}
