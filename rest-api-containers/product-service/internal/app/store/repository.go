package store

import "product-service/internal/app/model"

// ProductRepository представляет собой интерфейс для взаимодействия с хранилищем продуктов.
type ProductRepository interface {
	Create(*model.Product) error            // Создание нового продукта
	FindOne(string) (*model.Product, error) // Поиск продукта по идентификатору
	FindAll() ([]*model.Product, error)     // Поиск всех продуктов
	Delete(string) (*model.Product, error)  // Удаление продукта по идентификатору
}
