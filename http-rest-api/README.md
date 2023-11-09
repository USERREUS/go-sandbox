# Задание

Разработать RestAPI сервисов (Product Service, Order Service, Inventory Service, NotificationService)

## Описание

- ProductService      (PostgreSQL) 
- OrderService        (MongoDB)
- InventoryService    (PostgreSQL)
- NotificationService (MongoDB)


### Все сервисы поддерживают следующие HTTP маршруты:

- POST:   localhost:<8080>/\<service>      (create)
- GET:    localhost:<8080>/\<service>      (getAll)
- GET:    localhost:<8080>/\<service>/{id} (getOne)
- DELETE: localhost:<8080>/\<service>/{id} (Delete)
