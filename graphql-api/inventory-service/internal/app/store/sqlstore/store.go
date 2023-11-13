package sqlstore

import (
	"database/sql"
	"inventory/internal/app/store"
)

// Store представляет собой реализацию интерфейса store.Store для работы с базой данных SQL.
type Store struct {
	db         *sql.DB     // Подключение к базе данных
	repository *Repository // Репозиторий для взаимодействия с данными
}

// New создает новый экземпляр Store с указанным подключением к базе данных.
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// Repository возвращает репозиторий для взаимодействия с данными.
func (s *Store) Repository() store.Repository {
	if s.repository != nil {
		return s.repository
	}

	// Инициализация репозитория с использованием текущего экземпляра Store
	s.repository = &Repository{
		store: s,
	}

	return s.repository
}
