package models;

import (
	"log"
)

type WatchQueue struct {
	WatchQueueId int `db:"WatchQueueId"`
	UserId int `db:"UserId"`
	VideoId int `db:"VideoId"`
	DateAdded string `db:"DateAdded"`
	DateModified string `db:"DateModified"`
}

func GetWatchQueue(userId int, videoId int) (*WatchQueue, bool) {
	db := GetDBContext()

	sql := `SELECT WatchQueueId, UserId, VideoId, DateAdded, DateModified
	FROM WatchQueue
	WHERE UserId=? AND VideoId=?;`

	stmt, err := db.Preparex(sql)
	if err != nil {
		log.Fatal(err)
		return nil, false
	}

	var watchQueue *WatchQueue = new(WatchQueue)
	err = stmt.Get(watchQueue, userId, videoId)
	if err != nil {
		log.Println(err)
		log.Println(watchQueue)
		return nil, false
	}

	return watchQueue, true
}

func UpdateWatchQueue(watchQueue *WatchQueue) (*WatchQueue, bool) {
	var sql string
	if watchQueue.WatchQueueId <= 0 {
		sql = `INSERT INTO WatchQueue
			(VideoId, UserId, DateAdded, DateModified)
			VALUES
			(?,?,datetime(),datetime());`;
	} else {
		sql = `DELETE FROM WatchQueue
			WHERE UserId=? 
			AND VideoId=?;`;
	}

	db := GetDBContext()
	stmt, err := db.Preparex(sql)
	if err != nil {
		return watchQueue, false
	}

	if watchQueue.WatchQueueId <= 0 {
		result, err := stmt.Exec(watchQueue.VideoId, watchQueue.UserId)
		if err != nil {
			log.Print(err)
			return nil, false
		}
		insertId, err := result.LastInsertId()
		if err != nil {
			return watchQueue, false
		}
		watchQueue.WatchQueueId = int(insertId)
	} else {
		_, err := stmt.Exec(watchQueue.UserId, watchQueue.VideoId)
		if err != nil {
			log.Print(err)
			return watchQueue, false
		}
	}
	return watchQueue, true
}
