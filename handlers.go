package main


import (
	"fmt"
	"database/sql"
	"log"
	_ "github.com/glebarez/go-sqlite"
	"time"
	"strings"
)

const DATE_FORMAT  = "2006-01-02"
// Guarantee the Database is configurated for the tasks to be stored
func InitDB(filepath string) *sql.DB {

	//IF the .db file does not exist, the sql.Open function will create it automatically
	db, err_database_load := sql.Open("sqlite", filepath)
	if err_database_load != nil {
		log.Fatal("Error connecting to the database, make sure the sqlite database is configured correctly: ", err_database_load)
	}
	sql_script := `
	CREATE TABLE IF NOT EXISTS todo_items (
		todo_id INTEGER PRIMARY KEY,
		text_todo TEXT NOT NULL,
		complete INTEGER NOT NULL DEFAULT 0,
		due_date TEXT
	);`
	_, err_table_create := db.Exec(sql_script);
	if err_table_create != nil {
		log.Fatalf("Erro ao criar tabela: %q: %s\n", err_table_create, sql_script)
	}
	return db
}

// Function to add tasks into the Database
func AddTask(db *sql.DB, name string, due_date string) {
	sql_stat, err := db.Prepare("INSERT INTO tasks(name, date) VALUES(?, ?);")

	if err != nil {
		log.Fatal(err)
	}
	defer sql_stat.Close()

	_, err_insert := sql_stat.Exec(name, due_date)
	if err_insert != nil {
		log.Fatal(err_insert)
	}

	fmt.Printf("Task '%s' added succesfully!\n", name)
	if due_date != "" {
		fmt.Printf("\tDue Date: %s", due_date)
	}
}

func ListTasks(db *sql.DB) {
	rows, err_list := db.Query("SELECT todo_id, text_todo, complete, due_date FROM todo_items;")
	if err_list != nil {
		log.Fatal(err_list)
	}
	defer rows.Close()

}
func RemoveTask(db *sql.DB, task_ids []string) {
	var query string = ` DELETE FROM todo_items WHERE todo_id IN (`;
	query += strings.Join(task_ids, ", ")
	query += ")"
	_, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Task removed successfully")
	}

}

func CleanTasks(db *sql.DB) {
	var curr_date string = time.Now().Format(DATE_FORMAT)
	rows, err := db.Query("DELETE FROM todo_items WHERE due_date < ?", curr_date)
	if err != nil {
		log.Fatal(err)
	}
	var deleted_rows uint8 = 0
	for rows.Next() {
		deleted_rows++
	}
	fmt.Printf("%v deleted tasks\n", deleted_rows)
}
