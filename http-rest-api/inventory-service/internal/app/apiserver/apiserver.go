package apiserver

import (
	"database/sql"
	"inventory/internal/app/store/sqlstore"
	"net/http"

	_ "github.com/lib/pq"
)

// Start запускает API-сервер с указанной конфигурацией.
func Start(config *Config) error {
	// Инициализация подключения к базе данных
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()

	// Инициализация хранилища с использованием SQL-хранилища
	store := sqlstore.New(db)

	// Создание сервера с переданным хранилищем
	srv := newServer(store)

	// Запуск сервера на указанном адресе прослушивания
	return http.ListenAndServe(config.BindAddr, srv)
}

// newDB создает новое подключение к базе данных и выполняет необходимые инициализационные шаги.
func newDB(dbURL string) (*sql.DB, error) {
	// Открытие соединения с PostgreSQL базой данных
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	// Проверка подключения к базе данных
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Создание таблицы "inventory", если она не существует
	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS inventory" +
			"(" +
			"id VARCHAR NOT NULL, " +
			"name VARCHAR NOT NULL, " +
			"count DECIMAL NOT NULL, " +
			"cost DECIMAL NOT NULL" +
			")",
	); err != nil {
		return nil, err
	}

	return db, nil
}
