package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var connectionInstance *sql.DB = nil

func GetInstance() *sql.DB {
	var err error

	connectionInstance, err = sql.Open("mysql", "root:123@tcp(localhost:3306)/rest?parseTime=True")
	if err != nil {
		panic(err)
	}

	return connectionInstance
}
