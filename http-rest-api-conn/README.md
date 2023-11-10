# Задание

- Order во время формирования заказа должен организовать запрос к сервису Inventory, используя синхронное сообщение (rest api, grpc). 
- Сервис Order после формирования заказа должен посылать асинхронное сообщение (Kafka, RabbitMQ) сервису Notification.

# Описание

## Взимодействие

- Во время формирования заказа Order направляет запрос к Inventory, используя синхронное REST API сообщение. Проверяет наличие, количество и цену товара на складе. В случае успешной проверки добавляет заказ.
- После формирования заказа сервис Order формирует асинхронное сообщение (RabbitMQ) сервису Notification.

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
-   POST: /product 	-- Create(Product) -> Product
-    GET: /product 	-- GetAll() 	   -> map\[Code\]Poduct
-    GET: /product/{id} -- GetOne(Code)    -> Product
- DELETE: /product/{id} -- Delete(Code)    -> Product
```


## InventoryService

### localhost:8083

### PostgreSQL

### Inventory
```
{ 
	"code"  : string,
	"name" 	: string,
	"count" : int,
	"cost"  : int
}
```
	
### HTTP endpoints
```
-   POST: /inventory      -- Create(Inventory)   -> Inventory
-    PUT: /inventory      -- Update(Inventory)   -> Inventory
-    GET: /inventory      -- GetAll() 	       -> map\[ProductCode\]Inventory
-    GET: /inventory/{id} -- GetOne(ProductCode) -> Inventory
- DELETE: /inventory/{id} -- Delete(ProductCode) -> Inventory
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
-    GET: /notification       -- GetAll() 	      -> \[\]Notification
-    GET: /notification/{msg} -- GetMany(MsgType)     -> \[\]Notification
- DELETE: /notification/{msg} -- Delete(MsgType)      -> \[\]Notification
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
			"code"  : string,
			"name"  : string,
			"count" : int,
			"cost"  : int
		},...
	]
}
```
	
### HTTP endpoints
```
-   POST: /order      -- Create(\[\]ProductItem) -> \[\]ProductItem
-    GET: /order      -- GetAll() 	       -> map\[code\]Order
-    GET: /order/{id} -- GetOne(Code)          -> Order
- DELETE: /order/{id} -- Delete(Code)          -> Order
```
