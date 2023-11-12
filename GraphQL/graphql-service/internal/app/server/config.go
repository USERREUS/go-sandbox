package server

// Config представляет собой структуру конфигурации для API-сервера.
type Config struct {
	Handler     string `toml:"graphql_handler"` // Путь к обработчику GraphQL
	Port        string `toml:"port"`            // Порт, на котором работает API-сервер
	OrderAddr   string `toml:"order_addr"`      // Адрес сервиса заказов
	ProductAddr string `toml:"product_addr"`    // Адрес сервиса продуктов
}

// NewConfig создает новый экземпляр конфигурации с значениями по умолчанию.
func NewConfig() *Config {
	return &Config{
		Handler:     "/graphql",                      // Путь к обработчику GraphQL по умолчанию
		Port:        ":8080",                         // Порт по умолчанию
		OrderAddr:   "http://localhost:8081/order",   // Адрес сервиса заказов по умолчанию
		ProductAddr: "http://localhost:8084/product", // Адрес сервиса продуктов по умолчанию
	}
}
