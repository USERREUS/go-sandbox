package mongostore

import (
	"context"
	"notification-service/internal/app/store"

	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	context    context.Context
	collection *mongo.Collection
	repository *Repository
}

func New(context context.Context, collection *mongo.Collection) *Store {
	return &Store{
		context:    context,
		collection: collection,
	}
}

func (s *Store) Repository() store.Repository {
	if s.repository != nil {
		return s.repository
	}

	s.repository = &Repository{
		store: s,
	}

	return s.repository
}
