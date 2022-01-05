package renders

import "net/http"

func writeContentType(writer http.ResponseWriter, t string) {
	writer.Header().Del("Content-Type")
	writer.Header().Set("Content-Type", t)
}
