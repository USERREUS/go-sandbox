// Пакет main является точкой входа для выполнения программы на языке Go.
package main

// Импортирование необходимых пакетов для создания GraphQL-сервера.
import (
	"graphql-servive/internal/server" // Подключение локального пакета server с определением GraphQL-схемы.
	"log"                             // Пакет для логирования.
	"net/http"                        // Пакет для работы с HTTP-протоколом.

	"github.com/graphql-go/graphql" // Библиотека для работы с GraphQL.
	"github.com/graphql-go/handler" // Обработчик HTTP-запросов GraphQL.
)

// Главная функция main - точка входа программы на языке Go.
func main() {
	// Создание GraphQL-схемы на основе Query и Mutation из пакета server.
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    server.Query,
		Mutation: server.Mutation,
	})
	if err != nil {
		log.Fatal(err.Error()) // В случае ошибки создания схемы, программа завершается с выводом сообщения об ошибке.
	}

	// Создание нового обработчика GraphQL с конфигурацией.
	h := handler.New(&handler.Config{
		Schema:   &schema, // Указание схемы GraphQL для обработчика.
		Pretty:   true,    // Параметр Pretty для красивого форматирования ответов GraphQL.
		GraphiQL: true,    // Включение GraphiQL - веб-интерфейса для отладки GraphQL-запросов.
	})

	// Регистрация обработчика GraphQL по указанному адресу.
	http.Handle(server.Graphql_addr, h)

	// Вывод сообщения о запуске сервера GraphQL с указанием адреса и порта.
	log.Printf("GraphQL server is running on http://localhost%s/graphql\n", server.Graphql_port)

	// Запуск HTTP-сервера для обслуживания GraphQL-запросов.
	http.ListenAndServe(server.Graphql_port, nil)
}
