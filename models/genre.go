package models;

import (
	_ "github.com/mattn/go-sqlite3"
)

type Genre struct {
	GenreId int `db:"GenreId"`
	Description string `db:"Description"`
}

func GetRandomGenres() []Genre {
	db := GetDBContext()
	sql := `SELECT * 
		FROM Genre
		ORDER BY RANDOM()
		LIMIT 10`

	genres := []Genre{}
	err := db.Select(&genres, sql)
	if err != nil {
		panic(err)
	}

	db.Close()

	return genres
}
