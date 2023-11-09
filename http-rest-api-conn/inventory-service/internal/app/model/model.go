package model

// Model представляет собой структуру данных для хранения информации о записи инвентаря.
type Model struct {
	ID    string `json:"id,omitempty"` // Уникальный идентификатор записи (пустой, если не установлен)
	Name  string `json:"name"`         // Наименование
	Count int    `json:"count"`        // Количество
	Cost  int    `json:"cost"`         // Стоимость
}
