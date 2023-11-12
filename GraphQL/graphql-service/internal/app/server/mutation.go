package server

import (
	"errors"
	"graphql-service/internal/app/model"

	"github.com/graphql-go/graphql"
)

var Mutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"addOrder": &graphql.Field{
			Type: OrderType,
			Args: graphql.FieldConfigArgument{
				"orderInput": &graphql.ArgumentConfig{
					Type: OrderInputType,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				orderInput, ok := params.Args["orderInput"].(map[string]interface{})
				if !ok {
					return nil, errors.New("Invalid orderInput")
				}

				items, ok := orderInput["items"].([]interface{})
				if !ok {
					return nil, errors.New("Missing or invalid 'items' field in orderInput")
				}

				orderItems := make([]*model.OrderItem, len(items))

				for i, item := range items {
					itemMap, ok := item.(map[string]interface{})
					if !ok {
						return nil, errors.New("Invalid order item format")
					}

					orderItems[i] = &model.OrderItem{
						Code:  itemMap["code"].(string),
						Name:  itemMap["name"].(string),
						Count: int(itemMap["count"].(int)),
						Cost:  int(itemMap["cost"].(int)),
					}
				}

				order, err := AddOrder(orderItems)
				if err != nil {
					return nil, err
				}

				return order, nil
			},
		},
	},
})
