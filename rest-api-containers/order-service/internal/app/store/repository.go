package store

import "order-service/internal/app/model"

// Repository представляет интерфейс для работы с данными в хранилище.
type Repository interface {
	Create([]*model.ProductItem) (*model.Model, error) // Create создает новый заказ в хранилище и возвращает его уникальный идентификатор.
	FindOne(string) (*model.Model, error)              // FindOne находит заказ в хранилище по его уникальному идентификатору.
	FindAll() ([]*model.Model, error)                  // FindAll возвращает все заказы из хранилища.
	Delete(string) (*model.Model, error)               // Delete удаляет заказ из хранилища по его уникальному идентификатору.
}
