package apiserver

import (
	"context"
	"log"
	"net/http"
	"order-service/internal/app/store/mongostore"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client     *mongo.Client
	Context    context.Context
	Collection *mongo.Collection
}

// Start запускает API-сервер с указанной конфигурацией.
func Start(config *Config) error {
	mongodb, err := ConnectToMongoDB("restapi_dev", "order", config.DatabaseURL)
	if err != nil {
		return err
	}
	defer func() {
		if err := mongodb.Client.Disconnect(mongodb.Context); err != nil {
			log.Fatal("Error disconnecting from MongoDB:", err)
		}
	}()

	// Проверка подключения
	err = mongodb.Client.Ping(mongodb.Context, nil)
	if err != nil {
		log.Fatal("Unable to ping MongoDB:", err)
	}

	log.Println("Successfully connected to MongoDB")

	// Создание хранилища MongoDB
	store := mongostore.New(mongodb.Context, mongodb.Collection)

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

func ConnectToMongoDB(databaseName, collectionName, databaseURL string) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(databaseURL)
	ctx := context.Background()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	database := client.Database(databaseName)
	collection := database.Collection(collectionName)

	return &MongoDB{
		Client:     client,
		Context:    ctx,
		Collection: collection,
	}, nil
}
