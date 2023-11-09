package store

import "http-rest-api/internal/app/model"

type ProductRepository interface {
	Create(*model.Product) error
	FindOne(string) (*model.Product, error)
	Delete(string) (*model.Product, error)
	FindAll() (map[string]*model.Product, error)
}
