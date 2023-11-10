package model

// Product представляет собой структуру данных для описания продукта.
type Product struct {
	Code        string `json:"code"`        // Уникальный код продукта
	Name        string `json:"name"`        // Наименование продукта
	Weight      int    `json:"weight"`      // Вес продукта
	Description string `json:"description"` // Описание продукта
}
