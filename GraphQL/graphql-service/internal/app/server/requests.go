package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"graphql-service/internal/app/model"
	"net/http"
)

// GetOrders получает список заказов с удаленного сервиса заказов.
func (s *Server) GetOrders() ([]*model.Order, error) {
	resp, err := http.Get(s.Config.OrderAddr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var orders []*model.Order
	err = json.NewDecoder(resp.Body).Decode(&orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOrder получает заказ по указанному коду с удаленного сервиса заказов.
func (s *Server) GetOrder(code string) (*model.Order, error) {
	url := fmt.Sprintf("%s/%s", s.Config.OrderAddr, code)
	resp, err := http.Get(url)
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

// GetProducts получает список продуктов с удаленного сервиса продуктов.
func (s *Server) GetProducts() ([]*model.Product, error) {
	resp, err := http.Get(s.Config.ProductAddr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var products []*model.Product
	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

// GetProduct получает продукт по указанному коду с удаленного сервиса продуктов.
func (s *Server) GetProduct(code string) (*model.Product, error) {
	url := fmt.Sprintf("%s/%s", s.Config.ProductAddr, code)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var product model.Product
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// AddOrder добавляет заказ с указанными элементами заказа на удаленный сервис заказов.
func (s *Server) AddOrder(items []*model.OrderItem) (*model.Order, error) {
	s.Logger.Infoln("AddOrder")

	jsonData, err := json.Marshal(&items)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(s.Config.OrderAddr, "application/json", bytes.NewBuffer(jsonData))
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
