package renders

import (
	"github.com/eatmoreapple/regia/internal"
	"net/http"
)

type XmlRender struct {
	Serializer internal.Serializer
}

func (x XmlRender) Render(writer http.ResponseWriter, data interface{}) error {
	writeContentType(writer, "text/xml;charset=utf-8")
	return x.Serializer.Encode(writer, data)
}
