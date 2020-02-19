package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func main() {

	godotenv.Load()

	conn := os.Getenv("MYSQL_DSN")

	fmt.Println(conn)
	
}
