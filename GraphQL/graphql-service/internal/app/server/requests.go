package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"graphql-service/internal/app/model"
	"log"
	"net/http"
)

func GetOrders() ([]*model.Order, error) {
	resp, err := http.Get(Order_addr)
	if err != nil {
		log.Println("Error fetching orders:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var orders []*model.Order
	err = json.NewDecoder(resp.Body).Decode(&orders)
	if err != nil {
		log.Println("Error decoding orders response:", err)
		return nil, err
	}

	return orders, nil
}

func GetOrder(code string) (*model.Order, error) {
	url := fmt.Sprintf("%s/%s", Order_addr, code)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching order", err)
		return nil, err
	}
	defer resp.Body.Close()
	var order model.Order
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		log.Println("Error decoding order response:", err)
		return nil, err
	}

	return &order, nil
}

func GetProducts() ([]*model.Product, error) {
	resp, err := http.Get(Product_addr)
	if err != nil {
		log.Println("Error fetching products:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var products []*model.Product
	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		log.Println("Error decoding products response:", err)
		return nil, err
	}

	return products, nil
}

func GetProduct(code string) (*model.Product, error) {
	url := fmt.Sprintf("%s/%s", Product_addr, code)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching product", err)
		return nil, err
	}
	defer resp.Body.Close()
	var product model.Product
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		log.Println("Error decoding product response:", err)
		return nil, err
	}

	return &product, nil
}

func AddOrder(items []*model.OrderItem) (*model.Order, error) {
	jsonData, err := json.Marshal(&items)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(Order_addr, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var order model.Order
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
