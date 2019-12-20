package database

import (
	"database/sql"
	"github.com/davyxu/golog"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var log = golog.New("database")

func init() {
	var err error
	db, err = sql.Open("mysql", "gram:yangshu88@tcp(127.0.0.1:3306)/gram_landlord")
	if err != nil {
		log.Errorln(err)
	}
}