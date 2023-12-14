package main

import (
	"log"
	"notification-service/internal/app/apiserver"
)

func main() {
	// Создание объекта конфигурации
	config := apiserver.NewConfig()
	// Запуск сервера с использованием конфигурации
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
