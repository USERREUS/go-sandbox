package apiserver

import (
	"context"
	"log"
	"net/http"
	"notification-service/internal/app/store/mongostore"

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
	collection := database.Collection("notification")

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

	// Создание сервера с переданным хранилищем
	srv := newServer(store, ch)

	go func() {
		err := srv.Listen()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	// Запуск сервера на указанном адресе прослушивания
	return http.ListenAndServe(config.BindAddr, srv)
}
