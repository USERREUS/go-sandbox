package mongostore

import (
	"context"
	"notification-service/internal/app/store"

	"go.mongodb.org/mongo-driver/mongo"
)

// Store представляет собой хранилище данных MongoDB.
type Store struct {
	context    context.Context
	collection *mongo.Collection
	repository *Repository
}

// New создает новый экземпляр хранилища данных MongoDB.
func New(context context.Context, collection *mongo.Collection) *Store {
	return &Store{
		context:    context,
		collection: collection,
	}
}

// Repository возвращает интерфейс Repository для взаимодействия с хранилищем данных MongoDB.
func (s *Store) Repository() store.Repository {
	if s.repository != nil {
		return s.repository
	}

	s.repository = &Repository{
		store: s,
	}

	return s.repository
}
