package model

// Model представляет собой структуру данных для хранения информации о записи инвентаря.
type Model struct {
	Code  string `json:"code"`  // Kод продукта
	Name  string `json:"name"`  // Наименование
	Count int    `json:"count"` // Количество
	Cost  int    `json:"cost"`  // Стоимость
}
