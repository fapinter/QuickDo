package main


import (
	"fmt"
	"database/sql"
	"log"
	_ "github.com/glebarez/go-sqlite"
	"time"
	"strings"
	"regexp"
)

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
func AddTask(db *sql.DB, tasks[]string){

	var sql_script string = "INSERT INTO tasks(text_todo, due_date) VALUES"
	var tasks_inserted []string

	for i := 0; i < len(tasks); i++ {
		var date, text_todo string
			
		text_todo = tasks[i]
		re := regexp.MustCompile(":(\d{4}-\d{2}-\d{2})")
		dates := re.FindAll(text_todo)

		if len(dates) > 0 {
			date = dates[len(dates) - 1]
			text_todo = strings.TrimSuffix(text_todo, date)
		} else {
			date = "NULL"
		}
		var temp string = "(" + text_todo + "," + date + ")"
		tasks_inserted.append(temp)
	}
	for task_inserted {

	}
	_, err_insert := db.Exec(sql_script, text_todo, due_date)
	if err_insert != nil {
		log.Fatal(err_insert)
	}

	fmt.Printf("Task '%s' added succesfully!\n", len(tasks))
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

const DATE_FORMAT  = "2006-01-02"
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
