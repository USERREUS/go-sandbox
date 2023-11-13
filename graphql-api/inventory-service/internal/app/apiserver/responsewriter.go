package apiserver

import "net/http"

// responseWriter является оберткой для http.ResponseWriter с дополнительным полем для отслеживания статус-кода ответа.
type responseWriter struct {
	http.ResponseWriter
	code int
}

// WriteHeader переопределяет метод WriteHeader внутри responseWriter для сохранения статус-кода.
func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
