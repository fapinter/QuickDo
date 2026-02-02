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
		complete TEXT NOT NULL DEFAULT "No",
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
		var date, text_todo, temp string

		text_todo = value
		re := regexp.MustCompile(`:(\d{4}-\d{2}-\d{2})`)
		var dates [][]string = re.FindAllStringSubmatch(text_todo, -1)
		if len(dates) > 0 {
			//Separates the whole string on [0] and the capturing group on [1]
			date = dates[len(dates) - 1][1]
			var trim string = dates[len(dates) -1][0]
			text_todo = strings.TrimSuffix(text_todo, trim)

			//Validates if the Date does exist
			_, err_parse := time.Parse(DATE_FORMAT, date)
			if err_parse != nil {
				fmt.Printf("Invalid date %s, date will set as Today, use update-date to modify it\n", date)
				date = time.Now().Format(DATE_FORMAT)
			}

		}
		temp = fmt.Sprintf("('%s', '%s')", text_todo, date)
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
		)
		err := rows.Scan(&task.id, &task.text, &task.complete, &task.due_date)
		if err != nil {
			log.Fatalln(err)
		}
		if task.due_date != ""{
			date_task, err_parse = time.Parse(DATE_FORMAT, task.due_date)
			if err_parse != nil {
				fmt.Printf("Date %s could not be parsed, use the following format: YYYY-MM-DD\n", task.due_date)
			}
			if time.Now().After(date_task){
				task.expired = "Yes"
			} else{
				task.expired = "No"
			}
		} else{
			task.expired = "No"
		}

		fmt.Fprintf(w, "%v\t%s\t%s\t%v\t%s\t\n", task.id, task.text, task.complete, task.due_date, task.expired)
	}
	w.Flush()
}


func UpdateTask(db *sql.DB, id int, column string, value string) {
	//Date validation
	if column == "due_date"{
		_, err_parse := time.Parse(DATE_FORMAT, value)
		if err_parse != nil {
			fmt.Println("Invalid date, please use the following format: YYYY-MM-DD")
			return
		}
	}
	var query string = fmt.Sprintf("UPDATE todo_items SET %s='%s' WHERE todo_id=%v", column, value, id)
	res, err_query := db.Exec(query)
	if err_query != nil {
		log.Fatalln(err_query)
	}else {
		rows_affected, err_rows := res.RowsAffected()
		if err_rows != nil {
			log.Fatalln(err_rows)
		} else if rows_affected > 0{
			fmt.Printf("Task %v updated successfully\n", id)
		} else{
			fmt.Printf("Task %v not updated, check if the value passed is valid\n", id)
		}
	}
}

func ManageCheck(db *sql.DB, id int, check_state string) {
	var query string = fmt.Sprintf("UPDATE todo_items SET complete='%s' WHERE todo_id=%v", check_state, id)
	res, err := db.Exec(query)
	if err != nil {
		log.Fatalln(err)
	} else{
		rows_affected, err_rows := res.RowsAffected()
		if err_rows != nil {
			log.Fatalln(err_rows)
		} else if rows_affected > 0 {
			fmt.Printf("Task %v state changed successfully\n", id)
		} else {
			fmt.Printf("Task %v state not changed, task might not exist\n", id)
		}
	}
}


func RemoveTask(db *sql.DB, task_ids []string) {
	var query string = ` DELETE FROM todo_items WHERE todo_id IN (`;
	query += strings.Join(task_ids, ", ")
	query += ")"
	res, err_query := db.Exec(query)
	if err_query != nil {
		log.Fatalln(err_query)
	} else {
		rows_affected, err_rows := res.RowsAffected()
		if err_rows != nil {
			log.Fatalln(err_rows)
		} else {
			fmt.Printf("%d Task(s) removed successfully\n", rows_affected)
		}
	}

}

func CleanTasks(db *sql.DB, mode string) {
	var (
		deleted_rows int64
		err_rowCount error
		query string = "DELETE FROM todo_items WHERE "
	)
	if mode == "expired" {
		var curr_date string = time.Now().Format(DATE_FORMAT)
		query += fmt.Sprintf("due_date < '%s' AND due_date != ''", curr_date)
	} else{
		query += "complete='Yes'"
	}
	result, err := db.Exec(query)
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
