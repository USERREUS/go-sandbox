package sqlstore

import (
	"database/sql"
	"product-service/internal/app/store"
)

// Store представляет собой структуру для работы с SQL-хранилищем и предоставляет методы для получения репозитория продуктов.
type Store struct {
	db                *sql.DB
	productRepository *ProductRepository
}

// New создает новый экземпляр Store с указанным объектом базы данных.
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// Product возвращает репозиторий продуктов для взаимодействия с продуктами в хранилище.
func (s *Store) Product() store.ProductRepository {
	if s.productRepository != nil {
		return s.productRepository
	}

	s.productRepository = &ProductRepository{
		store: s,
	}

	return s.productRepository
}
