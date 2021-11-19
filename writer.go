package regia

import (
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) Write(data []byte) (int, error) {
	if w.status != 0 {
		w.WriteHeader(w.status)
	}
	return w.ResponseWriter.Write(data)
}
