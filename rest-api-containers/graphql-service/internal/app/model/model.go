package model

// Product структура представляет информацию о продукте
type Product struct {
	Code        string `json:"code"`        // Уникальный код продукта
	Name        string `json:"name"`        // Название продукта
	Weight      int    `json:"weight"`      // Вес продукта
	Description string `json:"description"` // Описание продукта
}

// OrderItem структура представляет элемент заказа
type OrderItem struct {
	Code  string `json:"product_code"` // Код продукта в элементе заказа
	Name  string `json:"name"`         // Название элемента заказа
	Count int    `json:"count"`        // Количество продуктов в заказе
	Cost  int    `json:"cost"`         // Стоимость элемента заказа
}

// Order структура представляет заказ
type Order struct {
	OrderCode  string       `json:"order_code"` // Уникальный код заказа
	Date       string       `json:"date"`       // Дата создания заказа
	OrderItems []*OrderItem `json:"order_item"` // Список элементов заказа
}
