package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"text/tabwriter"
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

	if len(os.Args) == 1 {
		fmt.Println("No command passed, use help to see all commands available")
		return
	}
	var op string = os.Args[1]

	switch op {
	case "add":
		var tasks []string = os.Args[2:]
		AddTask(db, tasks)

	case "list", "ls":
		var (
			cap_ int
			min  bool = false
			err  error
		)
		if len(os.Args) == 3 {
			if os.Args[2] == "--min" {
				min = true
			} else {
				cap_, err = strconv.Atoi(os.Args[2])
			}
		}
		if len(os.Args) > 3 {
			if contains(os.Args[2:], "--min") {
				min = true
			}
			var array_idx int = findInt(os.Args[2:])
			if array_idx == -1 {
				cap_ = -1
			} else {
				cap_, err = strconv.Atoi(os.Args[array_idx+2])
			}
		}
		if err != nil {
			fmt.Printf(INVALID_VALUE_INTEGER, "cap")
		} else {
			ListTasks(db, cap_, min)
		}

	case "check", "ck", "uncheck", "uck":
		var id_tasks []string
		var check_state string
		if op == "check" || op == "ck" {
			check_state = "Yes"
		} else {
			check_state = "No"
		}
		if len(os.Args) >= 3 {
			for idx, value := range os.Args[2:] {
				_, err_parse := strconv.Atoi(os.Args[idx+2])
				if err_parse != nil {
					fmt.Printf("%s is not a valid ID, will be ignored\n", os.Args[idx+2])
				} else {
					id_tasks = append(id_tasks, value)
				}
			}
			if len(id_tasks) > 0 {
				ManageCheck(db, id_tasks, check_state)
			} else {
				fmt.Printf("No valid IDs were passed, no task was (un)checked")
			}
		} else {
			log.Fatalln("Task ID missing, no task was (un)checked")
		}

	case "update-date", "up-dt", "update-text", "up-txt":
		var column, value string
		if len(os.Args) >= 4 {
			if len(os.Args) > 4 {
				fmt.Println("Update statements accept only one ID at a time, other parameters will be ignored")
			}
			id_task, err_parse := strconv.Atoi(os.Args[2])
			if err_parse != nil {
				fmt.Printf(INVALID_VALUE_INTEGER, "ID")
				return
			}
			value = os.Args[3]
			if op == "update-date" || op == "up-dt" {
				column = "due_date"
			} else {
				column = "text_todo"
			}
			UpdateTask(db, id_task, column, value)

		} else {
			if op == "update-date" {
				log.Fatalln("Task ID and/or Date missing, no task was updated")
			} else {
				log.Fatalln("Task ID and/or Text missing, no task was updated")
			}
		}

	case "remove", "rm":
		if len(os.Args) >= 3 {
			//IDs are passed as strings to facilitate the SQL statements
			//But must be verified so a SQL syntax error doesn't occur
			var tasks []string
			for _, value := range os.Args[2:] {
				if _, err := strconv.Atoi(value); err != nil {
					fmt.Printf("Invalid value %s for ID, must be integer\n", value)
				} else {
					tasks = append(tasks, value)
				}
			}
			RemoveTask(db, tasks)
		} else {
			fmt.Println("Task(s) ID(s) missing, no task was removed")
		}

	case "cleanup-expired", "cl-exp", "cleanup-completed", "cl-com":
		if op == "cleanup-expired" || op == "cl-exp" {
			CleanTasks(db, "expired")
		} else {
			CleanTasks(db, "completed")
		}
	case "alias":
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "COMMAND\tALIASES\t")
		fmt.Fprintln(w, "-------\t---------\t")
		fmt.Fprintln(w, "list\tls\t")
		fmt.Fprintln(w, "check\tck\t")
		fmt.Fprintln(w, "uncheck\tuck\t")
		fmt.Fprintln(w, "update-date\tup-dt\t")
		fmt.Fprintln(w, "update-text\tup-txt\t")
		fmt.Fprintln(w, "remove\trm\t")
		fmt.Fprintln(w, "cleanup-expired\tcl-exp\t")
		fmt.Fprintln(w, "cleanup-completed\tcl-com\t")
		fmt.Fprintln(w, "help\t--help, -h\t")
		fmt.Fprintln(w, "version\t--version, -v\t")
		fmt.Fprintln(w) // Extra blank line for readability
		w.Flush()

	case "help", "--help", "-h":
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

		fmt.Fprintln(w, "COMMAND\tDESCRIPTION\t")
		fmt.Fprintln(w, "-------\t---------\t")
		fmt.Fprintln(w, "add <text_task> ...\tAdd tasks, a date can be passed using : as separator (Ex: add \"clean room:2026-01-10\")\t")
		fmt.Fprintln(w, "list <cap> --min \tLists all the tasks created. Cap limits tasks shown, --min only shows ID and Text (Ex: list 4 --min)\t")
		fmt.Fprintln(w, "check <id>\tSet a task or more as completed (Ex: check 2 3 4)\t")
		fmt.Fprintln(w, "uncheck <id>\tSet a task or more as not completed (Ex: uncheck 2 3 4)\t")
		fmt.Fprintln(w, "update-date <id> <date>\tUpdates the due date of an existing task (Ex: update-date 2 2026-02-22)\t")
		fmt.Fprintln(w, "update-text <id> <text>\tUpdates the text of an existing task (Ex: update-text 2 wash the dishes)\t")
		fmt.Fprintln(w, "remove <id> ...\tRemoves one or more existing tasks (Ex: remove 2 1 4 5)\t")
		fmt.Fprintln(w, "cleanup-expired\tRemove all the expired tasks\t")
		fmt.Fprintln(w, "cleanup-completed\tRemove all the completed tasks\t")
		fmt.Fprintln(w, "alias\tDisplays all aliases to commands\t")
		fmt.Fprintln(w, "help\tDisplays a list of all the commands\t")
		fmt.Fprintln(w, "version\tDisplays the current version of Quickdo\t")
		fmt.Fprintln(w) // Extra blank line for readability
		w.Flush()

	case "version", "--version", "-v":
		fmt.Printf("Quickdo version: %s\n", Version)

	default:
		fmt.Println("Operation not supported, use the 'help' command to view all the operations possible")
	}
}

func contains(tasks []string, x string) bool {
	for _, n := range tasks {
		if x == n {
			return true
		}
	}
	return false
}

// Search linearly for the first interger in the string array
// returns -1 if not found
func findInt(tasks []string) int {
	for idx, n := range tasks {
		_, err_parse := strconv.Atoi(n)
		if err_parse == nil {
			return idx
		}
	}
	return -1
}
