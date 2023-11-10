package apiserver

// Config представляет структуру для хранения конфигурационных параметров API-сервера.
type Config struct {
	BindAddr    string `toml:"bind_addr"`    // Адрес прослушивания сервера
	LogLevel    string `toml:"log_level"`    // Уровень логирования
	DatabaseURL string `toml:"database_url"` // URL для подключения к MongoDB
	RabbitMQURL string `toml:"rabbitmq_url"` // URL для подключения к RabbitMQ
}

// NewConfig создает новый экземпляр Config с значениями по умолчанию.
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}
