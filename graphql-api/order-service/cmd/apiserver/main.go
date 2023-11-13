package main

import (
	"flag"
	"log"
	"order-service/internal/app/apiserver"

	"github.com/BurntSushi/toml"
)

// Определение переменных конфигурации.
var (
	configPath string
)

// Функция инициализации, устанавливающая флаги командной строки.
func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "путь к файлу конфигурации")
}

// Основная функция программы.
func main() {
	flag.Parse()

	// Создание нового экземпляра конфигурации.
	config := apiserver.NewConfig()

	// Чтение конфигурации из файла и декодирование ее в структуру Config.
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	// Запуск сервера API с использованием конфигурации.
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
