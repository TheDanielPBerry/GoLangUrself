package models


import (
	"github.com/jmoiron/sqlx"
	"log"
)

var DB *sqlx.DB

func connectDB() {
	var err error
	dbDriver, err := Config("DB", "DRIVER")
	if err != nil {
		panic(err)
	}

	dbPath, err := Config("DB", "PATH")
	if err != nil {
		panic(err)
	}

	DB, err = sqlx.Open(dbDriver, dbPath)
	if err != nil {
		panic(err)
	}
}

func GetDBContext() *sqlx.DB {
	if DB == nil {
		connectDB()
	} else if err := DB.Ping(); err != nil {
		connectDB()
	}
	return DB
}

func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Print(err)
		}
		DB = nil
	}
}
