package mongostore

import (
	"context"
	"order-service/internal/app/store"

	"go.mongodb.org/mongo-driver/mongo"
)

// Store представляет собой реализацию интерфейса store.Store для работы с MongoDB.
type Store struct {
	context    context.Context   // Контекст приложения
	collection *mongo.Collection // Коллекция MongoDB
	repository *Repository       // Репозиторий для взаимодействия с данными
}

// New создает новый экземпляр Store с указанным контекстом и коллекцией MongoDB.
func New(context context.Context, collection *mongo.Collection) *Store {
	return &Store{
		context:    context,
		collection: collection,
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
