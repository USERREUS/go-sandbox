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

	conn, err := amqp.Dial(config.RabbitMQURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"order",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Создание сервера с переданным хранилищем
	srv := newServer(store, ch)

	// Запуск сервера на указанном адресе прослушивания
	return http.ListenAndServe(config.BindAddr, srv)
}
