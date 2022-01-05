package renders

import "net/http"

type Render interface {
	Render(writer http.ResponseWriter, data interface{}) error
}
