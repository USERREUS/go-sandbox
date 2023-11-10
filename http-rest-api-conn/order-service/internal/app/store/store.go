package store

// Store представляет интерфейс для взаимодействия с хранилищем данных.
type Store interface {
	// Repository возвращает объект Repository для работы с данными.
	Repository() Repository
}
