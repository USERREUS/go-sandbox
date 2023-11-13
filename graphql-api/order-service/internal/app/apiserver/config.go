package apiserver

// Config представляет конфигурацию для API-сервера.
type Config struct {
	BindAddr    string `toml:"bind_addr"`    // Адрес прослушивания сервера
	LogLevel    string `toml:"log_level"`    // Уровень логирования
	DatabaseURL string `toml:"database_url"` // URL для подключения к базе данных
	RabbitMQURL string `toml:"rabbitmq_url"` // URL для подключения к RabbitMQ
}

// NewConfig создает новый экземпляр конфигурации с значениями по умолчанию.
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080", // По умолчанию прослушивать на порту 8080
		LogLevel: "debug", // Уровень логирования по умолчанию - debug
	}
}
