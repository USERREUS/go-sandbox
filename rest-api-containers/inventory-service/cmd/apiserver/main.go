package main

import (
	"inventory/internal/app/apiserver"
	"log"
)

func main() {
	// Создание нового экземпляра структуры Config
	config := apiserver.NewConfig()
	// Запуск API-сервера с предоставленной конфигурацией
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
