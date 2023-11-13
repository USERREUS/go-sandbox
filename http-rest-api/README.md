# Задание

Разработать RestAPI сервисов (Product Service, Order Service, Inventory Service, NotificationService)

## Описание

- ProductService      (PostgreSQL) 
- OrderService        (MongoDB)
- InventoryService    (PostgreSQL)
- NotificationService (MongoDB)


### Все сервисы поддерживают следующие HTTP маршруты:
```
- POST:   localhost:<port>/<service>      (create)
- GET:    localhost:<port>/<service>      (getAll)
- GET:    localhost:<port>/<service>/{id} (getOne)
- DELETE: localhost:<port>/<service>/{id} (Delete)
```
