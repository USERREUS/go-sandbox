package model

// Model представляет собой структуру данных для хранения информации о заказе в MongoDB.
type Model struct {
	OrderCode    string         `json:"order_code" bson:"order_code"`     // Уникальный код заказа (пустой, если не установлен)
	Date         string         `json:"date" bson:"date"`                 // Дата создания заказа
	ProductItems []*ProductItem `json:"product_item" bson:"product_item"` // Список товаров в заказе
}

// ProductItem представляет собой структуру данных для хранения информации о товаре в заказе.
type ProductItem struct {
	ProductCode string `json:"product_code" bson:"product_code"` // Уникальный код товара
	Name        string `json:"name" bson:"name"`                 // Наименование товара
	Count       int    `json:"count" bson:"count"`               // Количество товара
	Cost        int    `json:"cost" bson:"cost"`                 // Стоимость товара
}
