package server

import (
	"errors"
	"graphql-service/internal/app/model"

	"github.com/graphql-go/graphql"
)

// Mutation метод возвращает объект GraphQL для мутаций.
func (s *Server) Mutation() *graphql.Object {
	// Инициализация типов заказа и входных данных заказа для использования в мутации.
	orderType := s.Types.Order
	orderInputType := s.Types.OrderInput

	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			// Поле addOrder представляет мутацию для добавления заказа.
			"addOrder": &graphql.Field{
				Type: orderType,
				// Определение аргументов, которые могут быть переданы в мутацию.
				Args: graphql.FieldConfigArgument{
					"orderInput": &graphql.ArgumentConfig{
						Type: orderInputType,
					},
				},
				// Resolve функция обрабатывает мутацию, принимает аргументы и возвращает результат мутации.
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Извлечение orderInput из параметров запроса.
					orderInput, ok := params.Args["orderInput"].(map[string]interface{})
					if !ok {
						// В случае отсутствия или неверных данных в orderInput возвращается ошибка.
						return nil, errors.New("Invalid orderInput")
					}

					// Извлечение элементов заказа из orderInput.
					items, ok := orderInput["items"].([]interface{})
					if !ok {
						// В случае отсутствия или неверного формата данных в items возвращается ошибка.
						return nil, errors.New("Missing or invalid 'items' field in orderInput")
					}

					// Создание слайса для хранения элементов заказа.
					orderItems := make([]*model.OrderItem, len(items))

					// Итерация по каждому элементу в items для создания структуры OrderItem.
					for i, item := range items {
						itemMap, ok := item.(map[string]interface{})
						if !ok {
							// В случае неверного формата данных в элементе заказа возвращается ошибка.
							return nil, errors.New("Invalid order item format")
						}

						// Создание экземпляра OrderItem и заполнение его данными из itemMap.
						orderItems[i] = &model.OrderItem{
							Code:  itemMap["code"].(string),
							Name:  itemMap["name"].(string),
							Count: int(itemMap["count"].(int)),
							Cost:  int(itemMap["cost"].(int)),
						}
					}

					// Вызов метода сервера AddOrder для добавления заказа с элементами заказа.
					order, err := s.AddOrder(orderItems)
					if err != nil {
						// В случае ошибки при добавлении заказа в лог выводится сообщение об ошибке.
						s.Logger.Errorln("AddOrder Error: ", err)
						return nil, err
					}

					// Возвращение созданного заказа в результате мутации.
					return order, nil
				},
			},
		},
	})
}
