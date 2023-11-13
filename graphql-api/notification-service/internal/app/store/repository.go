package store

import "notification-service/internal/app/model"

// Repository представляет интерфейс для работы с хранилищем уведомлений.
type Repository interface {
	Create(*model.Model) error               // Создает новую запись уведомления
	FindMany(string) ([]*model.Model, error) // Находит все записи уведомлений по типу сообщения
	FindAll() ([]*model.Model, error)        // Находит все записи уведомлений
	Delete(string) error                     // Удаляет запись уведомления по типу сообщения
}
