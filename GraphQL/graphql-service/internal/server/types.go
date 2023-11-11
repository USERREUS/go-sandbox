package server

import "github.com/graphql-go/graphql"

var OrderType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Order",
	Fields: graphql.Fields{
		"order_code": &graphql.Field{
			Type: graphql.String,
		},
		"date": &graphql.Field{
			Type: graphql.String,
		},
		"order_item": &graphql.Field{
			Type: graphql.NewList(OrderItemType),
		},
	},
})

var OrderItemType = graphql.NewObject(graphql.ObjectConfig{
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

var ProductType = graphql.NewObject(graphql.ObjectConfig{
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

var OrderInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "OrderInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"items": &graphql.InputObjectFieldConfig{
			Type: graphql.NewList(OrderItemInputType),
		},
	},
})

var OrderItemInputType = graphql.NewInputObject(graphql.InputObjectConfig{
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
