package server

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/sirupsen/logrus"
)

// Server представляет собой структуру API-сервера.
type Server struct {
	Logger *logrus.Logger // Логгер для записи логов сервера.
	Config *Config        // Конфигурация API-сервера.
	Types  *GraphQLTypes  // Типы GraphQL, определенные для сервера.
}

// NewServer создает новый экземпляр API-сервера с указанной конфигурацией и логгером.
func NewServer(config *Config, logger *logrus.Logger) *Server {
	return &Server{
		Logger: logger,
		Config: config,
		Types:  TypesInit(),
	}
}

// Start запускает API-сервер с указанной конфигурацией.
func (s *Server) Start() error {
	// Создание GraphQL-схемы на основе Query и Mutation из пакета server.
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    s.Query(),
		Mutation: s.Mutation(),
	})
	if err != nil {
		return err
	}

	// Создание нового обработчика GraphQL с конфигурацией.
	h := handler.New(&handler.Config{
		Schema:   &schema, // Указание схемы GraphQL для обработчика.
		Pretty:   true,    // Параметр Pretty для красивого форматирования ответов GraphQL.
		GraphiQL: true,    // Включение GraphiQL - веб-интерфейса для отладки GraphQL-запросов.
	})

	// Регистрация обработчика GraphQL по указанному адресу.
	http.Handle(s.Config.Handler, h)

	// Вывод сообщения о запуске сервера GraphQL с указанием адреса и порта.
	logrus.Infof("GraphQL server is running on http://localhost%s/graphql\n", s.Config.Port)

	// Запуск HTTP-сервера для обслуживания GraphQL-запросов.
	return http.ListenAndServe(s.Config.Port, nil)
}
