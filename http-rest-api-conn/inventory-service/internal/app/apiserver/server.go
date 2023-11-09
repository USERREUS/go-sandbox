package apiserver

import (
	"encoding/json"
	"inventory/internal/app/model"
	"inventory/internal/app/store"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// server представляет собой основной сервер API с необходимыми полями и методами.
type server struct {
	router *mux.Router    // Роутер для маршрутизации HTTP-запросов
	logger *logrus.Logger // Логгер для записи логов
	store  store.Store    // Хранилище данных
}

// newServer создает новый экземпляр сервера с заданным хранилищем и конфигурирует роутер.
func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	s.configureRouter()

	return s
}

// ServeHTTP реализует интерфейс http.Handler для обработки HTTP-запросов через роутер сервера.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// configureRouter конфигурирует роутер сервера, устанавливая обработчики для различных маршрутов.
func (s *server) configureRouter() {
	s.router.Use(s.logRequest)
	s.router.HandleFunc("/inventory", s.handleInventoryCreate()).Methods("POST")
	s.router.HandleFunc("/inventory", s.handleInventoryFindAll()).Methods("GET")
	s.router.HandleFunc("/inventory/{id}", s.handleInventoryFindOne()).Methods("GET")
	s.router.HandleFunc("/inventory/{id}", s.handleInventoryDelete()).Methods("DELETE")
}

// logRequest является промежуточным слоем для логирования информации о запросе.
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

// handleInventoryCreate обрабатывает запрос на создание записи инвентаря.
func (s *server) handleInventoryCreate() http.HandlerFunc {
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

// handleInventoryFindOne обрабатывает запрос на поиск записи инвентаря по идентификатору.
func (s *server) handleInventoryFindOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, _ := mux.Vars(r)["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		record, err := s.store.Repository().FindOne(id.String())
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, record)
	}
}

// handleInventoryDelete обрабатывает запрос на удаление записи инвентаря по идентификатору.
func (s *server) handleInventoryDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, _ := mux.Vars(r)["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		record, err := s.store.Repository().Delete(id.String())
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, record)
	}
}

// handleInventoryFindAll обрабатывает запрос на получение списка всех записей инвентаря.
func (s *server) handleInventoryFindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		records, err := s.store.Repository().FindAll()
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, records)
	}
}

// error форматирует и возвращает ошибку в формате JSON.
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

// respond устанавливает заголовки, код статуса и отправляет данные в формате JSON клиенту.
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
