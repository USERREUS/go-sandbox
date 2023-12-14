package apiserver

import (
	"fmt"
	"os"
)

// Config представляет конфигурацию API-сервера.
type Config struct {
	BindAddr    string // Адрес прослушивания сервера
	LogLevel    string // Уровень логирования
	DatabaseURL string // URL для подключения к базе данных
}

// NewConfig создает новый экземпляр Config с значениями по умолчанию.
func NewConfig() *Config {
	// Получение значений переменных среды
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Формирование строки подключения
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	return &Config{
		BindAddr:    ":8080",
		LogLevel:    "debug",
		DatabaseURL: connStr,
	}
}
