package store

import "inventory/internal/app/model"

// Repository представляет интерфейс для взаимодействия с данными в хранилище.
type Repository interface {
	Create(*model.Model) error                 // Создание новой записи в хранилище
	Update(*model.Model) error                 // Обновление существующей записи в хранилище
	FindOne(string) (*model.Model, error)      // Поиск записи по идентификатору
	FindAll() (map[string]*model.Model, error) // Получение списка всех записей в хранилище
	Delete(string) (*model.Model, error)       // Удаление записи по идентификатору
}
