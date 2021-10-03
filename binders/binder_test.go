package binders

import (
	"net/url"
	"testing"
)

type User struct {
	Name string `form:"name"`
	Age  int    `form:"age"`
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
