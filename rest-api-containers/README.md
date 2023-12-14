# Задание

- Разработать GraphQL сервис для организации взаимодействия с сервисами (OrderService, ProductService), объединив их RestAPI в единую схему

# Описание

## Взимодействие

- Во время формирования заказа Order направляет запрос к Inventory, используя синхронное REST API сообщение. Проверяет наличие, количество и цену товара на складе. В случае успешной проверки добавляет заказ.
- После формирования заказа сервис Order формирует асинхронное сообщение (RabbitMQ) сервису Notification.
- GraphQL сервис запрашивает данные у Product и формирует заказ.

## ProductService 

### localhost:8084

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

### localhost:8083

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

### localhost:8082

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

### localhost:8081

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

- Find all orders
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
- Find order by code
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
- Find all products
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
- Find product by code
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

- Create order with slice order items
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
