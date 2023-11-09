package apiserver

// Config представляет конфигурацию API-сервера.
type Config struct {
	BindAddr    string `toml:"bind_addr"`    // Адрес прослушивания сервера
	LogLevel    string `toml:"log_level"`    // Уровень логирования
	DatabaseURL string `toml:"database_url"` // URL для подключения к базе данных
}

// NewConfig создает новый экземпляр Config с значениями по умолчанию.
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080", // По умолчанию прослушивает на порту 8080
		LogLevel: "debug", // Уровень логирования по умолчанию - debug
	}
}
