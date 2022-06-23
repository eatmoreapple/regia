package regia

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func runRequest(B *testing.B, r *Engine, method, path string) {
	// create fake request
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}
	w := httptest.NewRecorder()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		r.ServeHTTP(w, req)
	}

}

func BenchmarkOneRoute(B *testing.B) {
	router := New()
	router.GET("/ping", func(c *Context) {})
	router.init()
	runRequest(B, router, "GET", "/ping")
}

func BenchmarkRenderString(B *testing.B) {
	router := New()
	router.GET("/string", func(c *Context) { c.String("%s-%s", "hello", "world") })
	router.init()
	runRequest(B, router, "GET", "/string")
}

func BenchmarkRenderJson(B *testing.B) {
	router := New()
	router.GET("/json", func(c *Context) { c.JSON(Map{"hello": "world"}) })
	router.init()
	runRequest(B, router, "GET", "/json")
}
