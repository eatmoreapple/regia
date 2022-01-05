package regia

import (
	"github.com/eatmoreapple/regia/binders"
	"github.com/eatmoreapple/regia/internal"
)

type Binder interface {
	Bind(context *Context, v interface{}) error
}

type QueryBinder struct{}

func (QueryBinder) Bind(context *Context, v interface{}) error {
	query := context.Query()
	return binders.BindForm(query, v)
}

type FormBinder struct{}

func (FormBinder) Bind(context *Context, v interface{}) error {
	form := context.Form()
	return binders.BindForm(form, v)
}

type MultipartFormBodyBinder struct{}

func (MultipartFormBodyBinder) Bind(context *Context, v interface{}) error {
	return binders.BindMultipartForm(context.Request.MultipartForm, v)
}

type JsonBodyBinder struct {
	Serializer internal.Serializer
}

func (j JsonBodyBinder) Bind(context *Context, v interface{}) error {
	return j.Serializer.Decode(context.Request.Body, v)
}

type XmlBodyBinder struct {
	Serializer internal.Serializer
}

func (j XmlBodyBinder) Bind(context *Context, v interface{}) error {
	return j.Serializer.Decode(context.Request.Body, v)
}
