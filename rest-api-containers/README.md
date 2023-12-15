# Задание

- Создать контейнеры системы с помощью Docker. Развернуть архитектурное решение с помощью (Docker Compose и Kubernates)

# Описание

## Взимодействие

- Во время формирования заказа Order направляет запрос к Inventory, используя синхронное REST API сообщение. Проверяет наличие, количество и цену товара на складе. В случае успешной проверки добавляет заказ.
- После формирования заказа сервис Order формирует асинхронное сообщение (RabbitMQ) сервису Notification.
- GraphQL сервис запрашивает данные у Product и формирует заказ.

## Контейнеризация

- Для каждого сервиса создан свой Dockerfile.
- Для всей системы создан docker-compose.yaml файл.
- Внутри системы сервисы доступны по имени и порту.
- Вне системы доступен только graphql сервис на localhost:8080/graphql.

## ProductService 

### product-service:8080

### PostgreSQL

### Product
```
{ 
	"code"        : string,
	"name" 	      : string,
	"weight"      : int,
	"description" : string
}
```
	
### HTTP endpoints
```
-   POST: /product 	  -- Create(Product) -> Product
-    GET: /product 	  -- GetAll() 	     -> []*Poduct
-    GET: /product/{code} -- GetOne(Code)    -> Product
- DELETE: /product/{code} -- Delete(Code)    -> Product
```


## InventoryService

### inventory-service:8080

### PostgreSQL

### Inventory
```
{ 
	"product_code" : string,
	"name" 	       : string,
	"count"        : int,
	"cost"         : int
}
```
	
### HTTP endpoints
```
-   POST: /inventory        -- Create(Inventory)   -> Inventory
-    PUT: /inventory        -- Update(Inventory)   -> Inventory
-    GET: /inventory        -- GetAll() 	   -> []*Inventory
-    GET: /inventory/{code} -- GetOne(ProductCode) -> Inventory
- DELETE: /inventory/{code} -- Delete(ProductCode) -> Inventory
```


## NotificationService

### notification-service:8080

### MongoDB

### RabbitMQ

### Notification
```
{ 
	"message_type" : string,
	"description"  : string,
	"date" 	       : string
}
```
	
### HTTP endpoints
```
-   POST: /notification       -- Create(Notification) -> Notification
-    GET: /notification       -- GetAll() 	      -> []*Notification
-    GET: /notification/{msg} -- GetMany(MsgType)     -> []*Notification
- DELETE: /notification/{msg} -- Delete(MsgType)      -> []*Notification
```


## OrderService

### order-service:8080

### MongoDB

### RabbitMQ

### Order
```
{ 
	"order_code"   : string,
	"date"         : string,
	"product_item" : [
		{
			"product_code"  : string,
			"name"  : string,
			"count" : int,
			"cost"  : int
		},...
	]
}
```
	
### HTTP endpoints
```
-   POST: /order        -- Create([]*ProductItem) -> Order
-    GET: /order        -- GetAll() 	          -> []*Order
-    GET: /order/{code} -- GetOne(OrderCode)      -> Order
- DELETE: /order/{code} -- Delete(OrderCode)      -> Order
```


## GraphQLService

### localhost:8080

### Query

- Получить все заказы
```
{
  findAllOrders {
    order_code
    date
    order_item {
      product_code
      name
      count
      cost
    }
  }
}
```
- Получить заказ по коду
```
{
  findOrder(code: "<code>") {
    order_code
    date
    order_item {
      product_code
      name
      count
      cost
    }
  }
}

```
- Получить все продукты
```
{
  findAllProducts {
    code
    name
    weight
    description
  }
}
```
- Получить продукт по коду
```
{
  findProduct(code: "<code>") {
    code
    name
    weight
    description
  }
}
```

### Mutation

- Создать заказ по списку продуктов
```
mutation {
  addOrder (
    orderInput: {
      items: [
        {
          code: "<code>", 
          name: "<name>", 
          count: <count>, 
          cost:  <cost>
        },...
      ]}) 
  {
    order_code
    date
    order_item {
      product_code
      name
      count
      cost
    }
  }
}
```
