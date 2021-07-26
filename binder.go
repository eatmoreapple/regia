package regia

import "github.com/eatmoreapple/binder"

type Binder interface {
	Bind(request *Context, v interface{}) error
}

type QueryBinder struct{}

func (q QueryBinder) Bind(context *Context, v interface{}) error {
	return binder.BindForm(context.Query(), v)
}

type FormBinder struct{}

func (q FormBinder) Bind(context *Context, v interface{}) error {
	return binder.BindForm(context.Form(), v)
}

type JsonBodyBinder struct{}

func (q JsonBodyBinder) Bind(context *Context, v interface{}) error {
	return binder.JsonBodyBinder.Bind(context.Request.Body, v)
}

type XmlBodyBinder struct{}

func (x XmlBodyBinder) Bind(context *Context, v interface{}) error {
	return binder.XmlBodyBinder.Bind(context.Request.Body, v)
}

func AddCustomBindFormMethod(name string, method binder.CustomBindMethod) error {
	return binder.AddCustomBindFormMethod(name, method)
}

type MultipartFormBinder struct{}

func (m MultipartFormBinder) Bind(context *Context, v interface{}) error {
	return binder.BindMultipartForm(context.Request.MultipartForm, v)
}

var (
	queryBinder         = QueryBinder{}
	formBinder          = FormBinder{}
	jsonBinder          = JsonBodyBinder{}
	xmlBinder           = XmlBodyBinder{}
	multipartFormBinder = MultipartFormBinder{}
)
