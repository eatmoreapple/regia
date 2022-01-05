package renders

import (
	"fmt"
	"github.com/eatmoreapple/regia/internal"
	"net/http"
)

type StringRender struct {
	Format string
	Data   []interface{}
}

func (s StringRender) Render(writer http.ResponseWriter, v interface{}) (err error) {
	writeContentType(writer, "text/html;charset=utf-8")
	if len(s.Data) > 0 {
		_, err = fmt.Fprintf(writer, s.Format, s.Data...)
	} else {
		_, err = writer.Write(internal.StringToByte(s.Format))
	}
	return err
}
