package database

import (
	"database/sql"
	"main/helper"
)

func NewDB() *sql.DB {
	db, err := sql.Open("mysql", "root:kaliberjunior1@tcp(localhost:3306)/otp")
	helper.HandlePanic(err)

	return db
}
