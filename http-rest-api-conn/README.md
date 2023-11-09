# Задание

Order во время формирования заказа должен организовать запрос к сервису Inventory, используя синхронное сообщение (rest api, grpc). Сервис Order после формирования заказа должен посылать асинхронное сообщение (Kafka, RabbitMQ) сервису Notification.

## Описание

### Взимодействие

- Во время формирования заказа направляет запрос к , используя синхронное сообщение. Проверяет наличие, количество и цену товара на складе. В случае успешной проверки добавляет заказ.
- После формирования заказа сервис Order формирует асинхронное сообщение (RabbitMQ) сервису Notification.

### Хранение данных

- ProductService      (PostgreSQL) 
- OrderService        (MongoDB)
- InventoryService    (PostgreSQL)
- NotificationService (MongoDB)


### Все сервисы поддерживают следующие HTTP маршруты:

- POST:   localhost:<8080>/\<service>      (create)
- GET:    localhost:<8080>/\<service>      (getAll)
- GET:    localhost:<8080>/\<service>/{id} (getOne)
- DELETE: localhost:<8080>/\<service>/{id} (Delete)
