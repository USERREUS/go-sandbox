package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Строка подключения к PostgreSQL
	connStr := "user=postgres dbname=postgres sslmode=disable password=mysecretpassword host=postgres"

	// Открываем соединение с базой данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")

	// Здесь вы можете выполнять запросы и другие операции с базой данных
}
