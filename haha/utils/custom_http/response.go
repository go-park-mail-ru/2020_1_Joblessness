package custom_http

import "net/http"

type CustomResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func NewCustomResponseWriter(w http.ResponseWriter) *CustomResponseWriter {
	return &CustomResponseWriter{ResponseWriter: w}
}

func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
