package sqlstore

import (
	"database/sql"
	"http-rest-api/internal/app/store"
)

type Store struct {
	db                *sql.DB
	productRepository *ProductRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Product() store.ProductRepository {
	if s.productRepository != nil {
		return s.productRepository
	}

	s.productRepository = &ProductRepository{
		store: s,
	}

	return s.productRepository
}
