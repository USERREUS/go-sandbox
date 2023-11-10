package apiserver

import (
	"database/sql"
	"net/http"
	"product-service/internal/app/store/sqlstore"

	_ "github.com/lib/pq"
)

// Start запускает API-сервер с указанной конфигурацией.
func Start(config *Config) error {
	// Подключение к базе данных
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()

	// Создание хранилища SQL
	store := sqlstore.New(db)

	// Создание сервера с переданным хранилищем
	srv := newServer(store)

	// Запуск сервера на указанном адресе прослушивания
	return http.ListenAndServe(config.BindAddr, srv)
}

// newDB открывает соединение с базой данных PostgreSQL, выполняет ping и создает таблицу "products" при её отсутствии.
func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	// Проверка соединения с базой данных
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Создание таблицы "products" при её отсутствии
	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS products " +
			"(" +
			"code VARCHAR NOT NULL, " +
			"name VARCHAR NOT NULL, " +
			"weight INT NOT NULL, " +
			"description VARCHAR NOT NULL" +
			")",
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
