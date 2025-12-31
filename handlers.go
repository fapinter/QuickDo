package main


import (
	"fmt"
	"database/sql"
	"log"
	_ "github.com/glebarez/go-sqlite"
)


func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		log.Fatal("Error connecting to the database, make sure the sqlite database is configured correctly: ", err)
	}
	sql_script := `
	CREATE TABLE IF NOT EXISTS todo_items (
		todo_id INTEGER PRIMARY KEY,
		text_todo TEXT NOT NULL,
		complete INTEGER NOT NULL DEFAULT 0,
		due_date TEXT
	);`
	res, err := db.Exec(sql_script);
}
