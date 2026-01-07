package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

const DATE_FORMAT  = "2006-01-02"

type Task struct {
	id				int
	text  		string
	complete  string
	due_date  string
	expired   string
}

// Guarantee the Database is configurated for the tasks to be stored
func InitDB(filepath string) *sql.DB {

	//IF the .db file does not exist, the sql.Open function will create it automatically
	db, err_database_load := sql.Open("sqlite", filepath)
	if err_database_load != nil {
		log.Fatalln("Error connecting to the database, make sure the sqlite database is configured correctly: ", err_database_load)
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
		log.Fatalln(err_table_create)
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
		var dates [][]string = re.FindAllStringSubmatch(text_todo, -1)
		if len(dates) > 0 {
			//Separates the whole string on [0] and the capturing group on [1]
			date = dates[len(dates) - 1][1]
			var trim string = dates[len(dates) -1][0]
			text_todo = strings.TrimSuffix(text_todo, trim)
		}
		var temp string = fmt.Sprintf("('%s', '%s')", text_todo, date)
		tasks_inserted = append(tasks_inserted, temp)
	}
	sql_script += strings.Join(tasks_inserted, ",")
	_, err_insert := db.Exec(sql_script)
	if err_insert != nil {
		log.Fatalln(err_insert)
	}

	fmt.Printf("%d task(s) added successfully!\n", len(tasks))
}

func ListTasks(db *sql.DB, cap int) {
	var sql_script string = "SELECT todo_id, text_todo, complete, due_date FROM todo_items"
	if cap > 0 {
		sql_script += fmt.Sprintf(" LIMIT %d", cap)
	}
	rows, err_list := db.Query(sql_script)
	if err_list != nil {
		log.Fatalln(err_list)
	}
	defer rows.Close()
	
	//Tabwriter used to display the table aligned
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tText\tCompleted\tExpiration Date\tExpired\t")
	for rows.Next(){
		var (
			task Task
			date_task time.Time
			err_parse error
			completed bool
		)	
		err := rows.Scan(&task.id, &task.text, &completed, &task.due_date)
		if err != nil {
			log.Fatalln(err)
		}
		if task.due_date != ""{
			date_task, err_parse = time.Parse(DATE_FORMAT, task.due_date)
			if err_parse != nil {
				log.Fatalf("\nError parsing Task date: \n%s\n", err_parse)
			}
			if time.Now().After(date_task){
				task.expired = "Yes"
			} else{
				task.expired = "No"
			}	
		} else{
			task.expired = "No"
		}

		if completed{
			task.complete = "Yes"
		}	else{
			task.complete = "No"
		}
		fmt.Fprintf(w, "%v\t%s\t%s\t%v\t%s\t\n", task.id, task.text, task.complete, task.due_date, task.expired)
	}
	w.Flush()
}


func UpdateTask(db *sql.DB, id string, column string, value string) {
	var (
		query string = fmt.Sprintf("UPDATE todo_items SET %s=%s WHERE todo_id=%s", column, value, id)
	)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalln(err)
	}else {
		fmt.Printf("Task %s updated successfully", id)
	}
}
func ManageCheck(db *sql.DB, id string, check_state uint8) {
	var query string = "UPDATE todo_items SET checked=? WHERE todo_id=?"
	_, err := db.Query(query, check_state, id)
	if err != nil {
		log.Fatalf("Error (un)checking task: \n\t%s\n",err)
	} else{
		fmt.Printf("Task %v state changed successfully", id)
	}
}


func RemoveTask(db *sql.DB, task_ids []string) {
	var query string = ` DELETE FROM todo_items WHERE todo_id IN (`;
	query += strings.Join(task_ids, ", ")
	query += ")"
	_, err := db.Query(query)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Printf("%d Task(s) removed successfully", len(task_ids))
	}

}

func CleanTasks(db *sql.DB) {
	var (
		curr_date string = time.Now().Format(DATE_FORMAT)
		deleted_rows int64
		err_rowCount error
	)
	result, err := db.Exec("DELETE FROM todo_items WHERE due_date < ? AND due_date != ''", curr_date)
	if err != nil {
		log.Fatalln(err)
	}
	deleted_rows, err_rowCount = result.RowsAffected()
	if err_rowCount != nil {
		log.Fatalln(err_rowCount)
	} else{
		fmt.Printf("%d deleted tasks\n", deleted_rows)
	}
}
