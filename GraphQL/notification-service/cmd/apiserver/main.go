package main

import (
	"flag"
	"log"
	"notification-service/internal/app/apiserver"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	// Инициализация флагов командной строки
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "путь к файлу конфигурации")
}

func main() {
	// Разбор флагов командной строки
	flag.Parse()

	// Создание объекта конфигурации
	config := apiserver.NewConfig()

	// Загрузка конфигурации из файла
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	// Запуск сервера с использованием конфигурации
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
