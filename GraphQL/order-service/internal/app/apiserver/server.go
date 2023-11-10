package apiserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"order-service/internal/app/model"
	"order-service/internal/app/store"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// server - структура, представляющая собой сервер API.
type server struct {
	router *mux.Router    // Инстанс маршрутизатора Gorilla Mux.
	logger *logrus.Logger // Инстанс логгера Logrus.
	store  store.Store    // Интерфейс для взаимодействия с базой данных.
	ch     *amqp.Channel  // Канал RabbitMQ для отправки сообщений в очередь.
}

// newServer создает новый экземпляр сервера с переданным хранилищем и каналом RabbitMQ.
func newServer(store store.Store, ch *amqp.Channel) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
		ch:     ch,
	}

	// Настройка маршрутов.
	s.configureRouter()

	return s
}

// ServeHTTP реализует интерфейс http.Handler и позволяет серверу обрабатывать HTTP-запросы.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// configureRouter настраивает обработчики маршрутов для различных эндпоинтов API.
func (s *server) configureRouter() {
	s.router.Use(s.logRequest)
	s.router.HandleFunc("/order", s.handleOrderCreate()).Methods("POST")
	s.router.HandleFunc("/order", s.handleOrderFindAll()).Methods("GET")
	s.router.HandleFunc("/order/{code}", s.handleOrderFindOne()).Methods("GET")
	s.router.HandleFunc("/order/{code}", s.handleOrderDelete()).Methods("DELETE")
}

// logRequest - Middleware для логирования каждого входящего HTTP-запроса и его завершения.
func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логирование информации о запросе.
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		// Определение уровня логирования в зависимости от статусного кода ответа.
		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}

		// Логирование завершения запроса.
		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

// Enqueue отправляет сообщение в очередь RabbitMQ с информацией о заказе.
func (s *server) Enqueue(code string) error {
	body := fmt.Sprintf("order: %s", code)

	// Отправка сообщения в очередь RabbitMQ.
	err := s.ch.Publish(
		"",      // Обменник (пусто для очереди по умолчанию)
		"order", // Имя очереди
		false,   // Опубликовать ли сообщение, если нет потребителей
		true,    // Сообщение должно сохраняться при перезапуске сервера
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(body),
			DeliveryMode: amqp.Persistent, // persistent: гарантировать, что сообщение сохранится при перезапуске сервера
		})
	if err != nil {
		return err
	}

	log.Printf(" [x] Sent %s", body)

	return nil
}

// putInventoryProduct отправляет PUT-запрос к серверу инвентаря для обновления информации о продукте.
func (s *server) putInventoryProduct(port string, product *model.ProductItem) error {
	url := fmt.Sprintf("http://localhost:%s/inventory", port)

	// Преобразование данных в формат JSON.
	jsonData, err := json.Marshal(&product)
	if err != nil {
		return err
	}

	// Отправка PUT-запроса с использованием gorilla/mux.
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Установка заголовка Content-Type.
	req.Header.Set("Content-Type", "application/json")

	// Отправка запроса.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверка статуса ответа.
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Not found: %d", resp.StatusCode))
	}

	return nil
}

// getInventoryProductByCode отправляет GET-запрос к серверу инвентаря для получения информации о продукте по коду.
func (s *server) getInventoryProductByCode(port, code string) (*model.ProductItem, error) {
	url := fmt.Sprintf("http://localhost:%s/inventory/%s", port, code)

	// Отправка GET-запроса к серверу инвентаря.
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Парсинг JSON-ответа.
	var product model.ProductItem
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &model.ProductItem{
		ProductCode: product.ProductCode,
		Name:        product.Name,
		Count:       product.Count,
		Cost:        product.Cost,
	}, nil
}

// handleOrderCreate возвращает обработчик для создания заказа.
func (s *server) handleOrderCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var items []*model.ProductItem

		// Парсинг JSON-тела запроса в структуру заказа.
		if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		// Проверка наличия товаров в инвентаре и обновление информации о продуктах.
		// for _, item := range items {
		// 	prod, err := s.getInventoryProductByCode("8083", item.ProductCode)
		// 	if err != nil {
		// 		s.error(w, r, http.StatusNotFound, err)
		// 		return
		// 	}

		// 	if prod.Count < item.Count || prod.Cost != item.Cost {
		// 		s.error(w, r, http.StatusBadRequest, errors.New("Data error"))
		// 		return
		// 	}

		// 	prod.Count -= item.Count

		// 	err = s.putInventoryProduct("8083", prod)
		// 	if err != nil {
		// 		s.error(w, r, http.StatusNotFound, err)
		// 		return
		// 	}
		// }

		// Создание заказа в хранилище.
		res, err := s.store.Repository().Create(items)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		// // Отправка заказа в очередь RabbitMQ.
		// err = s.Enqueue(res.OrderCode)
		// if err != nil {
		// 	// Обработка ошибки (вывод в лог, но не прерывание выполнения).
		// 	log.Printf("Error enqueueing order: %v", err)
		// }

		// Отправка успешного ответа.
		s.respond(w, r, http.StatusCreated, res)
	}
}

// handleOrderFindOne возвращает обработчик для поиска одного заказа.
func (s *server) handleOrderFindOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code, _ := mux.Vars(r)["code"]
		record, err := s.store.Repository().FindOne(code)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, record)
	}
}

// handleOrderDelete возвращает обработчик для удаления заказа.
func (s *server) handleOrderDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code, _ := mux.Vars(r)["code"]
		record, err := s.store.Repository().Delete(code)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, record)
	}
}

// handleOrderFindAll возвращает обработчик для поиска всех заказов.
func (s *server) handleOrderFindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение всех заказов из хранилища.
		records, err := s.store.Repository().FindAll()
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		// Отправка успешного ответа с найденными заказами.
		s.respond(w, r, http.StatusOK, records)
	}
}

// error отправляет ошибку клиенту в формате JSON.
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	// Отправка ответа с кодом ошибки и сообщением об ошибке.
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

// respond отправляет успешный ответ клиенту в формате JSON.
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	// Установка заголовка Content-Type.
	w.Header().Set("Content-Type", "application/json")
	// Установка статусного кода ответа.
	w.WriteHeader(code)
	// Если есть данные, отправка их клиенту в формате JSON.
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
