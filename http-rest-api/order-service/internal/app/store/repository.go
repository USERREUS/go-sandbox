package store

import "order-service/internal/app/model"

type Repository interface {
	Create([]*model.ProductItem) error
	FindOne(string) (*model.Model, error)
	FindAll() (map[string]*model.Model, error)
	Delete(string) (*model.Model, error)
}
