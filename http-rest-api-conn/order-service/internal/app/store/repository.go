package store

import "order-service/internal/app/model"

// Repository представляет интерфейс для работы с данными в хранилище.
type Repository interface {
	// Create создает новый заказ в хранилище и возвращает его уникальный идентификатор.
	Create([]*model.ProductItem) (string, error)

	// FindOne находит заказ в хранилище по его уникальному идентификатору.
	FindOne(string) (*model.Model, error)

	// FindAll возвращает все заказы из хранилища.
	FindAll() (map[string]*model.Model, error)

	// Delete удаляет заказ из хранилища по его уникальному идентификатору.
	Delete(string) (*model.Model, error)
}
