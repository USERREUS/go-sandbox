package main

import (
	"log"
	"product-service/internal/app/apiserver"
)

// main является точкой входа в приложение.
func main() {
	// Создание экземпляра конфигурации и загрузка значений из файла.
	config := apiserver.NewConfig()
	// Запуск API-сервера с использованием конфигурации.
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
