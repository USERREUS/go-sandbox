package model

// Product структура представляет информацию о продукте
type Product struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Weight      int    `json:"weight"`
	Description string `json:"description"`
}

// OrderItem структура представляет элемент заказа
type OrderItem struct {
	Code  string `json:"product_code"`
	Name  string `json:"name"`
	Count int    `json:"count"`
	Cost  int    `json:"cost"`
}

// Order структура представляет заказ
type Order struct {
	OrderCode  string       `json:"order_code"`
	Date       string       `json:"date"`
	OrderItems []*OrderItem `json:"order_item"`
}
