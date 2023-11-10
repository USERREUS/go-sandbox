# Задание

Order во время формирования заказа должен организовать запрос к сервису Inventory, используя синхронное сообщение (rest api, grpc). Сервис Order после формирования заказа должен посылать асинхронное сообщение (Kafka, RabbitMQ) сервису Notification.

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

-   POST: /product 	-- Create(Product) -> Product
-    GET: /product 	-- GetAll() 	   -> map\[Code\]Poduct
-    GET: /product/{id} -- GetOne(Code)    -> Product
- DELETE: /product/{id} -- Delete(Code)    -> Product


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
-   POST: localhost:8083/inventory      -- Create(Inventory)   -> Inventory
-    PUT: localhost:8083/inventory      -- Update(Inventory)   -> Inventory
-    GET: localhost:8083/inventory      -- GetAll() 	       -> map\[ProductCode\]Inventory
-    GET: localhost:8083/inventory/{id} -- GetOne(ProductCode) -> Inventory
- DELETE: localhost:8083/inventory/{id} -- Delete(ProductCode) -> Inventory
```
