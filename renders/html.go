package renders

import (
	"html/template"
	"net/http"
)

type Template struct {
	*template.Template
}

func (t Template) Render(writer http.ResponseWriter, data interface{}) error {
	writeContentType(writer, "text/html; charset=utf-8")
	return t.Execute(writer, data)
}
