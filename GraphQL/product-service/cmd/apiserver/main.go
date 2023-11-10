package main

import (
	"flag"
	"log"
	"product-service/internal/app/apiserver"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

// init устанавливает значения флагов по умолчанию.
func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "путь к файлу конфигурации")
}

// main является точкой входа в приложение.
func main() {
	flag.Parse()

	// Создание экземпляра конфигурации и загрузка значений из файла.
	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	// Запуск API-сервера с использованием конфигурации.
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
