package store

// Store представляет интерфейс для работы с хранилищем данных.
type Store interface {
	Repository() Repository // Возвращает репозиторий для взаимодействия с данными
}
