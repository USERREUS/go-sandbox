package apiserver

import "net/http"

// responseWriter представляет структуру, расширяющую http.ResponseWriter для отслеживания статуса ответа.
type responseWriter struct {
	http.ResponseWriter
	code int
}

// WriteHeader переопределяет метод WriteHeader стандартного интерфейса http.ResponseWriter,
// чтобы сохранять код статуса ответа.
func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
