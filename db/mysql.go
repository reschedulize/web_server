package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var MySQL *sqlx.DB

func ConnectMySQL() {
	if MySQL == nil {
		MySQL = sqlx.MustConnect("mysql", "root:password@tcp(127.0.0.1:3306)/reschedulize")
	}
}
