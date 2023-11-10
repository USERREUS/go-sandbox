package apiserver

import "net/http"

// responseWriter представляет собой расширение http.ResponseWriter для отслеживания кода ответа.
type responseWriter struct {
	http.ResponseWriter // Встроенный интерфейс http.ResponseWriter
	code                int
}

// WriteHeader переопределяет метод WriteHeader для сохранения кода ответа.
func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
