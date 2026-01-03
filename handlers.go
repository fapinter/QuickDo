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
		log.Fatalf("Erro ao criar tabela: %s: %s\n", err_table_create, sql_script)
	}
	return db
}

// Function to add tasks into the Database
func AddTask(db *sql.DB, tasks[]string){
	var sql_script string = "INSERT INTO todo_items(text_todo, due_date) VALUES"
	var tasks_inserted []string
	for _, value := range tasks {
		var date, text_todo string
			
		text_todo = value
		re := regexp.MustCompile(`:(\d{4}-\d{2}-\d{2})`)
		var dates []string = re.FindAllString(text_todo, -1)
		//Makes sure to get the last one if a Date is found
		if len(dates) > 0 {
			date = dates[len(dates) - 1]
			text_todo = strings.TrimSuffix(text_todo, date)
		}
		var temp string = fmt.Sprintf("('%s', '%s')", text_todo, date)
		tasks_inserted = append(tasks_inserted, temp)
	}
	sql_script += strings.Join(tasks_inserted, ",")
	_, err_insert := db.Exec(sql_script)
	if err_insert != nil {
		log.Fatal(err_insert)
	}

	fmt.Printf("Task '%d' added successfully!\n", len(tasks))
}

func ListTasks(db *sql.DB, cap int) {
	var sql_script string = "SELECT todo_id, text_todo, complete, due_date FROM todo_items"
	if cap > 0 {
		sql_script += fmt.Sprintf(" LIMIT %d", cap)
	}
	rows, err_list := db.Query(sql_script)
	if err_list != nil {
		log.Fatal(err_list)
	}
	defer rows.Close()
	
	type Task struct {
		id				int
		text  		string
		completed bool
		due_date  string
		expired   bool
	}
	
}

func UpdateTask(db *sql.DB, id int, text string, date string) {
	var (
		query string = "UPDATE todo_items SET "
		err error
	)
	if text != "" && date != "" {
		query += "text_todo=?, due_date=? WHERE todo_id=?"
		_, err = db.Exec(query, text, date, id)
	} else if text != "" {
		query += "text_todo=? WHERE todo_id=?"
		_, err = db.Query(query, text, id)
	}else{
		query += "due_date=? WHERE todo_id=?"	
		_, err = db.Query(query, date, id)
	}
	if err != nil {
		log.Fatal(err)
	}else {
		fmt.Println("Task updated successfully")
	}
}


func RemoveTask(db *sql.DB, task_ids []string) {
	var query string = ` DELETE FROM todo_items WHERE todo_id IN (`;
	query += strings.Join(task_ids, ", ")
	query += ")"
	_, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%d Task(s) removed successfully", len(task_ids))
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
	fmt.Printf("%d deleted tasks\n", deleted_rows)
}
