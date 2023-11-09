package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func test() {
	// Устанавливаем соединение с сервером RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Создаем канал
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Объявляем очередь
	q, err := ch.QueueDeclare(
		"hello", // Имя очереди
		false,   // Постоянная ли очередь
		false,   // Удалить ли очередь, когда отписывается последний потребитель
		false,   // Продолжительность жизни сообщения
		false,   // Очередь автоматически удаляется при отсутствии потребителей
		nil,     // Аргументы
	)
	failOnError(err, "Failed to declare a queue")

	// Отправляем сообщение в очередь
	body := "Hello World!"
	err = ch.Publish(
		"",     // Обменник (пусто для очереди по умолчанию)
		q.Name, // Имя очереди
		false,  // Опубликовать ли сообщение, если нет потребителей
		false,  // Сообщение не должно сохраняться при перезапуске сервера
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)

	// Читаем сообщение из очереди
	msgs, err := ch.Consume(
		q.Name, // Имя очереди
		"",     // Имя потребителя (пусто для автоматического создания)
		true,   // Автоматическое подтверждение (ack)
		false,  // Исключение из очереди при неудачной обработке
		false,  // Не использовать канал с потоками
		false,  // Параметры конфликтов
		nil,    // Аргументы
	)
	failOnError(err, "Failed to register a consumer")

	for msg := range msgs {
		log.Printf(" [x] Received %s", msg.Body)
	}
}
func main() {
	go test()
	fmt.Println("Done")
}
