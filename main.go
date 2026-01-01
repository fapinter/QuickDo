package main


import (
	"fmt"
	
)

const DATABASE_PATH = "./todos.db"

func main() {
	InitDB(DATABASE_PATH)
	fmt.Print("Database initiated")
}
