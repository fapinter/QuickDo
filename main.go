package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const DATABASE_PATH = "./todos.db"
const CHECK_TASK uint8 = 1
const UNCHECK_TASK uint8 = 0
func main() {
	db := InitDB(DATABASE_PATH)

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
				log.Fatal(err)
			} else {
				ListTasks(db, cap_)
			}

		case "check", "uncheck":
			var id_task string
			var check_state uint8
			if op == "check"{
				check_state = CHECK_TASK
			}
			if len(os.Args) >= 3{
				if len(os.Args) > 3{
					fmt.Println("check function accepts only one ID at a time, other parameters will be ignored")
				}
				id_task = os.Args[2] 
				ManageCheck(db, id_task, CHECK_TASK)
			} else {
				log.Fatalln("Task ID missing, no task was checked")
			} 
			
		case "update-date":


		case "update-text":

		case "remove":
			var tasks []string = os.Args[2:]
			RemoveTask(db, tasks)

		case "cleanup":
			CleanTasks(db)

		case "export":
			fmt.Println("Work in progress")

		case "help":
			fmt.Println("Help not finished")
	}
}
