package main

import (
	"flag"
	"inventory/internal/app/apiserver"
	"log"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	// Определение флага командной строки для указания пути к файлу конфигурации
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "путь к файлу конфигурации")
}

func main() {
	// Разбор флагов командной строки
	flag.Parse()

	// Создание нового экземпляра структуры Config
	config := apiserver.NewConfig()

	// Декодирование конфигурации из указанного файла в структуру Config
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	// Запуск API-сервера с предоставленной конфигурацией
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
