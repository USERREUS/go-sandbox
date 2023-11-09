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

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
	ch     *amqp.Channel
}

func newServer(store store.Store, ch *amqp.Channel) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
		ch:     ch,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.logRequest)
	s.router.HandleFunc("/order", s.handleOrderCreate()).Methods("POST")
	s.router.HandleFunc("/order", s.handleOrderFindAll()).Methods("GET")
	s.router.HandleFunc("/order/{id}", s.handleOrderFindOne()).Methods("GET")
	s.router.HandleFunc("/order/{id}", s.handleOrderDelete()).Methods("DELETE")
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}
		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *server) Enqueue(id string) error {
	// Отправляем сообщение в очередь
	body := fmt.Sprintf("order: %s", id)
	err := s.ch.Publish(
		"",      // Обменник (пусто для очереди по умолчанию)
		"order", // Имя очереди
		false,   // Опубликовать ли сообщение, если нет потребителей
		false,   // Сообщение не должно сохраняться при перезапуске сервера
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		return err
	}

	log.Printf(" [x] Sent %s", body)

	return nil
}

func (s *server) putInventoryProduct(port string, product *model.ProductItem) error {
	// Формируем URL для запроса к серверу
	url := fmt.Sprintf("http://localhost:%s/inventory", port)

	// Преобразуем данные в формат JSON
	jsonData, err := json.Marshal(&product)
	if err != nil {
		return err
	}

	// Отправляем PUT-запрос с использованием gorilla/mux
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Устанавливаем заголовок Content-Type
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Not found: %d", resp.StatusCode))
	}

	return nil
}

func (s *server) getInventoryProductByCode(port, code string) (*model.ProductItem, error) {
	// Формируем URL для запроса к серверу
	url := fmt.Sprintf("http://localhost:%s/inventory/%s", port, code)

	// Отправляем GET-запрос к серверу
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Парсим JSON-ответ
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

func (s *server) handleOrderCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var items []*model.ProductItem
		if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		for _, item := range items {
			prod, err := s.getInventoryProductByCode("80XX", item.ProductCode)
			if err != nil {
				s.error(w, r, http.StatusNotFound, err)
				return
			}

			if prod.Count < item.Count || prod.Cost != item.Cost {
				s.error(w, r, http.StatusBadRequest, errors.New("Data error"))
				return
			}

			prod.Count -= item.Count

			err = s.putInventoryProduct("8080", prod)
			if err != nil {
				s.error(w, r, http.StatusNotFound, err)
				return
			}
		}

		id, err := s.store.Repository().Create(items)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		err = s.Enqueue(id)
		if err != nil {
			return
		}

		s.respond(w, r, http.StatusCreated, items)
	}
}

func (s *server) handleOrderFindOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, _ := mux.Vars(r)["id"]
		record, err := s.store.Repository().FindOne(idStr)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusFound, record)
	}
}

func (s *server) handleOrderDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, _ := mux.Vars(r)["id"]
		record, err := s.store.Repository().Delete(idStr)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusFound, record)
	}
}

func (s *server) handleOrderFindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		records, err := s.store.Repository().FindAll()
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusFound, records)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
