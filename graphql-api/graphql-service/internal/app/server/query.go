package server

import (
	"errors"

	"github.com/graphql-go/graphql"
)

// Query метод возвращает объект GraphQL для запросов.
func (s *Server) Query() *graphql.Object {
	// Инициализация типов продукта и заказа для использования в запросах.
	productType := s.Types.Product
	orderType := s.Types.Order

	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			// Поле findOrder представляет запрос для поиска заказа по коду.
			"findOrder": &graphql.Field{
				Type: orderType,
				// Определение аргументов, которые могут быть переданы в запросе.
				Args: graphql.FieldConfigArgument{
					"code": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				// Resolve функция обрабатывает запрос, извлекает код заказа из аргументов и возвращает соответствующий заказ.
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					code, ok := params.Args["code"].(string)
					if !ok {
						// В случае отсутствия или неверных данных в коде заказа возвращается ошибка.
						return nil, errors.New("Missing or invalid orderID")
					}

					// Вызов метода сервера GetOrder для получения заказа по коду.
					record, err := s.GetOrder(code)
					if err != nil {
						// В случае ошибки при выполнении запроса заказа в лог выводится сообщение об ошибке.
						s.Logger.Errorln("GetOrder Error:", err)
						return nil, err
					}

					// Возвращение найденного заказа в результате запроса.
					return record, nil
				},
			},
			// Поле findProduct представляет запрос для поиска продукта по коду.
			"findProduct": &graphql.Field{
				Type: productType,
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
					// Вызов метода сервера GetProduct для получения продукта по коду.
					record, err := s.GetProduct(code)
					if err != nil {
						s.Logger.Errorln("GetProduct Error:", err)
						return nil, err
					}

					// Возвращение найденного продукта в результате запроса.
					return record, nil
				},
			},
			// Поле findAllProducts представляет запрос для получения списка всех продуктов.
			"findAllProducts": &graphql.Field{
				Type: graphql.NewList(productType),
				Resolve: func(_ graphql.ResolveParams) (interface{}, error) {
					// Вызов метода сервера GetProducts для получения списка всех продуктов.
					records, err := s.GetProducts()
					if err != nil {
						s.Logger.Errorln("GetProducts Error:", err)
						return nil, err
					}

					// Возвращение списка всех продуктов в результате запроса.
					return records, nil
				},
			},
			// Поле findAllOrders представляет запрос для получения списка всех заказов.
			"findAllOrders": &graphql.Field{
				Type: graphql.NewList(orderType),
				Resolve: func(_ graphql.ResolveParams) (interface{}, error) {
					// Вызов метода сервера GetOrders для получения списка всех заказов.
					records, err := s.GetOrders()
					if err != nil {
						s.Logger.Error("GetOrders Error:", err)
						return nil, err
					}

					// Возвращение списка всех заказов в результате запроса.
					return records, nil
				},
			},
		},
	})
}
