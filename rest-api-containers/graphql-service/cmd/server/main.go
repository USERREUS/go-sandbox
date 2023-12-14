// Пакет main является точкой входа для выполнения программы на языке Go.
package main

// Импортирование необходимых пакетов для создания GraphQL-сервера.
import (
	// Подключение локального пакета server с определением GraphQL-схемы.

	"graphql-service/internal/app/server"

	// Пакет для логирования.

	"github.com/sirupsen/logrus"
	// Библиотека для работы с GraphQL.
	// Обработчик HTTP-запросов GraphQL.
)

// main является точкой входа в приложение.
func main() {
	// Инициализация логгера.
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	// Создание экземпляра конфигурации и загрузка значений из файла.
	config := server.NewConfig()
	// Создание экземпляра GraphQL-сервера с использованием конфигурации и логгера.
	srv := server.NewServer(config, logger)

	// Запуск API-сервера с использованием конфигурации.
	if err := srv.Start(); err != nil {
		// Обработка ошибки при запуске сервера.
		logger.Fatalln(err)
	}
}
