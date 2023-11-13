package store

// Store представляет собой интерфейс, предоставляющий метод Product для получения репозиторий продуктов.
type Store interface {
	Product() ProductRepository
}
