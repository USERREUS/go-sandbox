package apiserver

import (
	"context"
	"net/http"
	"order-service/internal/app/store/mongostore"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Start запускает API-сервер с указанной конфигурацией.
func Start(config *Config) error {
	// Настройка параметров подключения к MongoDB
	clientOptions := options.Client().ApplyURI(config.DatabaseURL)
	ctx := context.Background()
	defer ctx.Done()

	// Подключение к MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	// Выбор базы данных и коллекции в MongoDB
	database := client.Database("restapi_dev")
	collection := database.Collection("order")

	// Создание хранилища MongoDB
	store := mongostore.New(ctx, collection)

	// Подключение к RabbitMQ
	conn, err := amqp.Dial(config.RabbitMQURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Создание канала RabbitMQ
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Объявление очереди RabbitMQ для заказов
	_, err = ch.QueueDeclare(
		"order", // Имя очереди
		true,    // durable: Устойчивая (переживает перезапуск брокера)
		false,   // autoDelete: Автоматически удалять при отсутствии подключений
		false,   // exclusive: Использование только текущим подключением
		false,   // noWait: Не ждать подтверждения от брокера
		nil,     // Аргументы (nil означает использование значений по умолчанию)
	)

	if err != nil {
		return err
	}

	// Создание сервера с переданным хранилищем и каналом RabbitMQ
	srv := newServer(store, ch)

	// Запуск сервера на указанном адресе прослушивания
	return http.ListenAndServe(config.BindAddr, srv)
}
