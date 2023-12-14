package main

import (
	"log"
	"order-service/internal/app/apiserver"
)

// Основная функция программы.
func main() {
	// Создание нового экземпляра конфигурации.
	config := apiserver.NewConfig()
	// Запуск сервера API с использованием конфигурации.
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
