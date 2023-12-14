package server

import "fmt"

// Config представляет собой структуру конфигурации для API-сервера.
type Config struct {
	Handler     string // Путь к обработчику GraphQL
	Port        string // Порт, на котором работает API-сервер
	OrderAddr   string // Адрес сервиса заказов
	ProductAddr string // Адрес сервиса продуктов
}

// NewConfig создает новый экземпляр конфигурации с значениями по умолчанию.
func NewConfig() *Config {
	productServiceHost := "product-service"
	productServicePort := "8080"

	// Формирование URL для подключения к другому сервису
	productServiceURL := fmt.Sprintf("http://%s:%s/product", productServiceHost, productServicePort)

	return &Config{
		Handler:     "/graphql",                    // Путь к обработчику GraphQL по умолчанию
		Port:        ":8080",                       // Порт по умолчанию
		OrderAddr:   "http://localhost:8081/order", // Адрес сервиса заказов по умолчанию
		ProductAddr: productServiceURL,             // Адрес сервиса продуктов по умолчанию
	}
}
