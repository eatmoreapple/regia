package binders

import (
	"net/http"
	"net/url"
	"testing"
)

type User struct {
	Name string `form:"name" header:"name"`
	Age  int    `form:"age" header:"age"`
}

func TestBindForm(t *testing.T) {
	var user User
	var form = url.Values{"name": {"eatmoreapple"}, "age": {"18"}}
	if err := BindForm(form, &user); err != nil {
		t.Error(err)
		return
	}
	t.Log(user)
}

func TestHeader(t *testing.T) {
	var user User
	var form = http.Header{"name": {"eatmoreapple"}, "age": {"18"}}
	binder := DefaultHeaderBinder
	if err := binder.BindForm(url.Values(form), &user); err != nil {
		t.Error(err)
		return
	}
	t.Log(user)
}
