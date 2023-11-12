package server

// Config представляет собой структуру конфигурации для API-сервера.
type Config struct {
	BindAddr string `toml:"bind_addr"` // Адрес прослушивания сервера
}

// NewConfig создает новый экземпляр конфигурации с значениями по умолчанию.
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
	}
}
