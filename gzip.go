package regia

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	gzipWriter *gzip.Writer
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	return w.gzipWriter.Write(b)
}

func Gzip(level int) HandleFunc {
	return func(c *Context) {
		// Check if the client can accept gzip compression
		if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			return
		}
		writer, err := gzip.NewWriterLevel(c.ResponseWriter, level)
		if err != nil {
			c.Logger.Error(err)
			return
		}

		c.SetHeader("Content-Encoding", "gzip")
		c.SetHeader("Vary", "Accept-Encoding")

		// Wrap the response writer with the gzip writer
		w := &gzipWriter{c.ResponseWriter, writer}
		c.ResponseWriter = w
		defer func() {
			c.SetHeader("Content-Length", "0")
			if err = writer.Close(); err != nil {
				c.Logger.Error(err)
			}
		}()
		c.Next()
	}
}
