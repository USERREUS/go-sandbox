package server

import "github.com/graphql-go/graphql"

// GraphQLTypes структура, содержащая определения типов GraphQL для сервера.
type GraphQLTypes struct {
	Order          *graphql.Object      // Тип объекта для заказа
	OrderItem      *graphql.Object      // Тип объекта для элемента заказа
	Product        *graphql.Object      // Тип объекта для продукта
	OrderInput     *graphql.InputObject // Тип входного объекта для заказа
	OrderInputItem *graphql.InputObject // Тип входного объекта для элемента заказа
}

// TypesInit функция, инициализирующая и возвращающая экземпляр GraphQLTypes с определенными типами GraphQL.
func TypesInit() *GraphQLTypes {
	return &GraphQLTypes{
		OrderItem:      orderItemType,
		Order:          orderType,
		Product:        productType,
		OrderInput:     orderInputType,
		OrderInputItem: orderItemInputType,
	}
}

// Определения типов GraphQL для заказа, элемента заказа и продукта:

var orderType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Order",
	Fields: graphql.Fields{
		"order_code": &graphql.Field{
			Type: graphql.String,
		},
		"date": &graphql.Field{
			Type: graphql.String,
		},
		"order_item": &graphql.Field{
			Type: graphql.NewList(orderItemType),
		},
	},
})

var orderItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "OrderItem",
	Fields: graphql.Fields{
		"product_code": &graphql.Field{
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

var orderInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "OrderInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"items": &graphql.InputObjectFieldConfig{
			Type: graphql.NewList(orderItemInputType),
		},
	},
})

var orderItemInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "OrderItemInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"code": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"count": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
		"cost": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
	},
})
