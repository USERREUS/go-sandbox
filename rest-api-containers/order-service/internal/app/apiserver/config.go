package apiserver

import (
	"fmt"
	"os"
)

// Config представляет конфигурацию для API-сервера.
type Config struct {
	BindAddr    string // Адрес прослушивания сервера
	LogLevel    string // Уровень логирования
	DatabaseURL string // URL для подключения к базе данных
	RabbitMQURL string // URL для подключения к RabbitMQ
}

// NewConfig создает новый экземпляр конфигурации с значениями по умолчанию.
func NewConfig() *Config {

	// Получение значений переменных окружения
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")
	rabbitMQHost := os.Getenv("RABBITMQ_HOST")
	rabbitMQPort := os.Getenv("RABBITMQ_PORT")

	// Пример использования значений в коде
	mongoConnectionString := fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort)
	rabbitMQConnectionString := fmt.Sprintf("amqp://%s:%s", rabbitMQHost, rabbitMQPort)

	return &Config{
		BindAddr:    ":8080",
		LogLevel:    "debug",
		DatabaseURL: mongoConnectionString,
		RabbitMQURL: rabbitMQConnectionString,
	}
}
