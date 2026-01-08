package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"
	"path/filepath"
)

const DATABASE_PATH = "todos.db"
const INVALID_VALUE_INTEGER = "Invalid value for %s, must be an integer\n"

var Version = "development"

func getDBPath() string {
		user_path, _ := os.UserConfigDir()
		app_dir := filepath.Join(user_path, "quickdo")
		os.MkdirAll(app_dir, 0755)
		return filepath.Join(app_dir, DATABASE_PATH)
}

func main() {
	filepath := getDBPath()
	db := InitDB(filepath)

	if len(os.Args) == 1{
		fmt.Println("No command passed, use help to see all commands available")
		return
	}
 	var op string = os.Args[1]

	switch op {
		case "add":
			var tasks []string = os.Args[2:]
			AddTask(db, tasks)

		case "list":
			var (
				cap_ int
				err error
			)
			if len(os.Args) >= 3 {
				cap_, err = strconv.Atoi(os.Args[2])
			}
			if err != nil {
				fmt.Printf(INVALID_VALUE_INTEGER, "cap")
			} else {
				ListTasks(db, cap_)
			}

		case "check", "uncheck":
			var id_task int
			var check_state string
			var err_parse error
			if op == "check"{
				check_state = "Yes"
			} else {
				check_state = "No"
			}
			if len(os.Args) >= 3{
				if len(os.Args) > 3{
					fmt.Println("(un)check function accepts only one ID at a time, other parameters will be ignored")
				}
				id_task, err_parse = strconv.Atoi(os.Args[2]) 
				if err_parse != nil {
					fmt.Printf(INVALID_VALUE_INTEGER, "id")
				}
				
				ManageCheck(db, id_task, check_state)
			} else {
				log.Fatalln("Task ID missing, no task was (un)checked")
			} 
			
		case "update-date", "update-text":
			var column, value string
			if len(os.Args) >= 4{
				if len(os.Args) > 4 {
					fmt.Println("Update statements accept only one ID at a time, other parameters will be ignored")
				}
				id_task, err_parse := strconv.Atoi(os.Args[2])
				if err_parse != nil{
					fmt.Printf(INVALID_VALUE_INTEGER, "ID")
					return
				}
				value = os.Args[3]
				if op == "update-date"{
					column = "due_date"
				} else {
					column = "text_todo"
				}
				UpdateTask(db, id_task, column, value)

			} else {
				log.Fatalln("Task ID and/or Date missing, no task was updated")
			}


		case "remove":
			if len(os.Args) >= 3{
				//IDs are passed as strings to facilitate the SQL statements
				//But must be verified so a SQL syntax error doesn't occur
				var tasks []string
				for _, value := range os.Args[2:] {
					if _, err := strconv.Atoi(value); err != nil {
						fmt.Printf("Invalid value %s for ID, must be integer\n", value)
					} else{
						tasks = append(tasks, value)	
					}
				}
				RemoveTask(db, tasks)
			} else {
				fmt.Println("Task(s) ID(s) missing, no task was removed")
			}

		case "cleanup-expired", "cleanup-completed":
			if op == "cleanup-expired" {
				CleanTasks(db, "expired")
			} else {
				CleanTasks(db, "completed")
			}


		case "export":
			fmt.Println("Work in progress")

		case "help":
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

				fmt.Fprintln(w, "COMMAND\tDESCRIPTION\t")
				fmt.Fprintln(w, "-------\t---------\t")
				fmt.Fprintln(w, "add <text_task> ...\tAdd tasks, a date can be passed using : as separator (Ex: add \"clean room:2026-01-10\")\t")
				fmt.Fprintln(w, "list <cap>\tLists all the tasks created. Cap can be passed to limit how many tasks are shown (Ex: list 4)\t")
				fmt.Fprintln(w, "check <id>\tSet a task as completed (Ex: check 2)\t")
				fmt.Fprintln(w, "uncheck <id>\tSet a task as not completed (Ex: uncheck 2)\t")
				fmt.Fprintln(w, "update-date <id> <date>\tUpdates the due date of an existing task (Ex: update-date 2 2026-02-22)\t")
				fmt.Fprintln(w, "update-text <id> <text>\tUpdates the text of an existing task (Ex: update-text 2 wash the dishes)\t")
				fmt.Fprintln(w, "remove <id> ...\tRemoves one or more existing tasks (Ex: remove 2 1 4 5)\t")
				fmt.Fprintln(w, "cleanup-expired\tRemove all the expired tasks\t")
				fmt.Fprintln(w, "cleanup-completed\tRemove all the completed tasks\t")
				fmt.Fprintln(w, "export\tExports the tasks to a Todo List on MyNotes (still working on it)\t")
				fmt.Fprintln(w, "help\tDisplays a list of all the commands\t")
				fmt.Fprintln(w) // Linha em branco extra para legibilidade
				w.Flush()

		case "--version", "-v", "version":
			fmt.Printf("Quickdo version: %s\n", Version)

		default:
			fmt.Println("Operation not supported, use the help command to view all the operations possible")
	}
}
