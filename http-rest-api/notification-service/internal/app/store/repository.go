package store

import "notification-service/internal/app/model"

type Repository interface {
	Create(*model.Model) error
	FindMany(string) ([]*model.Model, error)
	FindAll() ([]*model.Model, error)
	Delete(string) error
}
