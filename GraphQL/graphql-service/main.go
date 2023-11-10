package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

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
	OrderCode    string       `json:"order_code"`
	Date         string       `json:"date"`
	ProductItems []*OrderItem `json:"product_item"`
}

var orderType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Order",
	Fields: graphql.Fields{
		"order_code": &graphql.Field{
			Type: graphql.String,
		},
		"date": &graphql.Field{
			Type: graphql.String,
		},
		"product_item": &graphql.Field{
			Type: graphql.NewList(orderItemType),
		},
	},
})

var orderItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "OrderItem",
	Fields: graphql.Fields{
		"code": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"count": &graphql.Field{
			Type: graphql.Int,
		},
		"cost": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

func getOrders() ([]*Order, error) {
	url := fmt.Sprintf("http://localhost:8081/order")

	// Отправка GET-запроса к серверу инвентаря.
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Парсинг JSON-ответа.
	var orders []*Order
	err = json.NewDecoder(resp.Body).Decode(&orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func getProducts() ([]*Product, error) {
	url := fmt.Sprintf("http://localhost:8084/product")

	// Отправка GET-запроса к серверу инвентаря.
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Парсинг JSON-ответа.
	var products []*Product
	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

var productType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Product",
	Fields: graphql.Fields{
		"code": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"weight": &graphql.Field{
			Type: graphql.Int,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var products = make(map[string]*Product)
var orders = make(map[string]*Order)

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"product": &graphql.Field{
			Type: productType,
			Args: graphql.FieldConfigArgument{
				"code": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				code, ok := params.Args["code"].(string)
				if !ok {
					return nil, nil
				}
				return products[code], nil
			},
		},
		"allProducts": &graphql.Field{
			Type: graphql.NewList(productType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// Генерация нового продукта со случайными данными
				records, err := getProducts()
				if err != nil {
					log.Fatal(err.Error())
				}

				// for _, record := range records {
				// 	products[record.Code] = record
				// }

				// // Формирование списка всех продуктов для возврата
				// var productList []*Product
				// for _, v := range products {
				// 	productList = append(productList, v)
				// }
				return records, nil
			},
		},
		"allOrders": &graphql.Field{
			Type: graphql.NewList(orderType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// Генерация нового продукта со случайными данными
				records, err := getOrders()
				if err != nil {
					log.Fatal(err.Error())
				}

				return records, nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
})

func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	var params graphql.Params
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := graphql.Do(params)
	json.NewEncoder(w).Encode(result)
}

func main() {
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", nil)
}
