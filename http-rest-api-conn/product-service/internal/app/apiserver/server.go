package apiserver

import (
	"encoding/json"
	"net/http"
	"product-service/internal/app/model"
	"product-service/internal/app/store"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// server представляет собой основной сервер API.
type server struct {
	router *mux.Router    // Роутер маршрутов
	logger *logrus.Logger // Логгер
	store  store.Store    // Хранилище данных
}

// newServer создает новый экземпляр сервера с указанным хранилищем данных.
func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(), // Создание нового роутера
		logger: logrus.New(),    // Создание нового логгера
		store:  store,           // Присвоение хранилища данных
	}

	s.configureRouter() // Настройка маршрутов

	return s
}

// ServeHTTP реализует интерфейс http.Handler, обрабатывая HTTP-запросы.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// configureRouter настраивает маршруты сервера.
func (s *server) configureRouter() {
	s.router.Use(s.logRequest)
	s.router.HandleFunc("/product", s.handleProductsCreate()).Methods("POST")
	s.router.HandleFunc("/product", s.handleProductsFindAll()).Methods("GET")
	s.router.HandleFunc("/product/{id}", s.handleProductsFindOne()).Methods("GET")
	s.router.HandleFunc("/product/{id}", s.handleProductsDelete()).Methods("DELETE")
}

// logRequest создает промежуточное ПО для записи логов каждого запроса.
func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
		})
		logger.Infof("начало обработки %s %s", r.Method, r.RequestURI)

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
			"завершено с кодом %d %s за %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

// handleProductsCreate обрабатывает запрос на создание нового продукта.
func (s *server) handleProductsCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &model.Product{}
		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Product().Create(p); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, p)
	}
}

// handleProductsFindOne обрабатывает запрос на поиск продукта по идентификатору.
func (s *server) handleProductsFindOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := mux.Vars(r)["id"]

		record, err := s.store.Product().FindOne(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, record)
	}
}

// handleProductsDelete обрабатывает запрос на удаление продукта по идентификатору.
func (s *server) handleProductsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := mux.Vars(r)["id"]

		record, err := s.store.Product().Delete(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, record)
	}
}

// handleProductsFindAll обрабатывает запрос на поиск всех продуктов.
func (s *server) handleProductsFindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		records, err := s.store.Product().FindAll()
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, records)
	}
}

// error отправляет HTTP-ответ с ошибкой.
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

// respond отправляет HTTP-ответ с указанным кодом состояния и данными.
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
