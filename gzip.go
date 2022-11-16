// Copyright 2022 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

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

// Write writes to the underlying gzip writer
func (w *gzipWriter) Write(b []byte) (int, error) {
	return w.gzipWriter.Write(b)
}

// Gzip is a middleware for gzip compression
// it will compress the response body if the client accepts gzip encoding
// param is the compression level, choose from gzip.BestSpeed to gzip.BestCompression
func Gzip(level int) HandleFunc {
	return func(c *Context) {
		// Check if the client can accept gzip compression
		if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			return
		}
		writer, err := gzip.NewWriterLevel(c.ResponseWriter, level)
		if err != nil {
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
				return
			}
		}()
		c.Next()
	}
}
