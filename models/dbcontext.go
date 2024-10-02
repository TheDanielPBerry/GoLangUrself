package models


import (
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB


func GetDBContext() *sqlx.DB {
	if DB == nil {
		var err error
		DB, err = sqlx.Open("sqlite3", "westflix.db")

		if err != nil {
			panic(err)
		}
	}
	return DB
}
