package apiserver

import (
	"context"
	"net/http"
	"order-service/internal/app/store/mongostore"

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

	// Создание сервера с переданным хранилищем
	srv := newServer(store)

	// Запуск сервера на указанном адресе прослушивания
	return http.ListenAndServe(config.BindAddr, srv)
}
