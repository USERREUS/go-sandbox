package apiserver

import (
	"encoding/json"
	"log"
	"net/http"
	"notification-service/internal/app/model"
	"notification-service/internal/app/store"
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

func (s *server) Listen() error {
	// Читаем сообщение из очереди
	msgs, err := s.ch.Consume(
		"order", // Имя очереди
		"",      // Имя потребителя (пусто для автоматического создания)
		true,    // Автоматическое подтверждение (ack)
		false,   // Исключение из очереди при неудачной обработке
		false,   // Не использовать канал с потоками
		false,   // Параметры конфликтов
		nil,     // Аргументы
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		log.Printf(" [x] Received %s", msg.Body)
		m := &model.Model{
			MsgType:     "order",
			Description: string(msg.Body),
		}
		s.store.Repository().Create(m)
	}

	return nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.logRequest)
	s.router.HandleFunc("/notification", s.handleNotificationCreate()).Methods("POST")
	s.router.HandleFunc("/notification", s.handleNotificationFindAll()).Methods("GET")
	s.router.HandleFunc("/notification/{msg}", s.handleNotificationFindMany()).Methods("GET")
	s.router.HandleFunc("/notification/{msg}", s.handleNotificationDelete()).Methods("DELETE")
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

func (s *server) handleNotificationCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &model.Model{}
		if err := json.NewDecoder(r.Body).Decode(m); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Repository().Create(m); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, m)
	}
}

func (s *server) handleNotificationFindMany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msgType, _ := mux.Vars(r)["msg"]
		records, err := s.store.Repository().FindMany(msgType)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusFound, records)
	}
}

func (s *server) handleNotificationDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msgType, _ := mux.Vars(r)["msg"]
		err := s.store.Repository().Delete(msgType)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusFound, "Delete success")
	}
}

func (s *server) handleNotificationFindAll() http.HandlerFunc {
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
