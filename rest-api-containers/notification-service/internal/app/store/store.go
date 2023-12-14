package store

// Store представляет интерфейс для работы с хранилищем.
type Store interface {
	Repository() Repository
}
