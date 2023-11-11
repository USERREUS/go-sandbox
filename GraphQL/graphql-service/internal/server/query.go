package server

import (
	"errors"
	"log"

	"github.com/graphql-go/graphql"
)

var Query = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"findOrder": &graphql.Field{
			Type: OrderType, // предположим, что у вас есть объект типа Order
			Args: graphql.FieldConfigArgument{
				"code": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				code, ok := params.Args["code"].(string)
				if !ok {
					return nil, errors.New("Missing or invalid orderID")
				}

				record, err := GetOrder(code)
				if err != nil {
					log.Println("Error:", err)
					return nil, err
				}

				return record, nil
			},
		},
		"findProduct": &graphql.Field{
			Type: ProductType,
			Args: graphql.FieldConfigArgument{
				"code": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				code, ok := params.Args["code"].(string)
				if !ok {
					return nil, errors.New("FindProduct Error")
				}
				record, err := GetProduct(code)
				if err != nil {
					log.Println("Error:", err)
					return nil, err
				}

				return record, nil
			},
		},
		"findAllProducts": &graphql.Field{
			Type: graphql.NewList(ProductType),
			Resolve: func(_ graphql.ResolveParams) (interface{}, error) {
				records, err := GetProducts()
				if err != nil {
					log.Println("Error:", err)
					return nil, err
				}

				return records, nil
			},
		},
		"findAllOrders": &graphql.Field{
			Type: graphql.NewList(OrderType),
			Resolve: func(_ graphql.ResolveParams) (interface{}, error) {
				records, err := GetOrders()
				if err != nil {
					log.Println("Error:", err)
					return nil, err
				}

				return records, nil
			},
		},
	},
})
