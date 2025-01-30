package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Ihu")
	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"locahost",
		5432,
		"postgres",
		"postgres",
		"dbname",
	)
	_, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Conecteido")
}
