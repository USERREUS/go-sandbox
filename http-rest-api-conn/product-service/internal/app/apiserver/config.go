package apiserver

// Config представляет собой структуру конфигурации для API-сервера.
type Config struct {
	BindAddr    string `toml:"bind_addr"`    // Адрес прослушивания сервера
	LogLevel    string `toml:"log_level"`    // Уровень логирования
	DatabaseURL string `toml:"database_url"` // URL базы данных
}

// NewConfig создает новый экземпляр конфигурации с значениями по умолчанию.
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}
