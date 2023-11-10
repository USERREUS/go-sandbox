package apiserver

import "net/http"

// responseWriter представляет собой структуру, расширяющую http.ResponseWriter.
type responseWriter struct {
	http.ResponseWriter
	code int
}

// WriteHeader переопределяет метод WriteHeader для отслеживания статусного кода ответа.
func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
