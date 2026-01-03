package main


import (
	"os"
	"strconv"
	"log"
)

const DATABASE_PATH = "./todos.db"

func main() {
	db := InitDB(DATABASE_PATH)

 	var op string = os.Args[1]

	switch op {
		case "add":
			var tasks []string = os.Args[2:]
			AddTask(db, tasks)

		case "list":
			if len(os.Args) == 3 {
				cap, err := strconv.Atoi(os.Args[2])
			}
			if err != nil {
				log.Fatal(err)
			} else {
				ListTasks(db, cap)
			}

		case "update":

		case "remove":

		case "cleanup":

		case "export":
	}
}
