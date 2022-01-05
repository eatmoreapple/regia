package renders

import "net/http"

type RedirectRender struct {
	RedirectURL string
	Code        int
	Request     *http.Request
}

func (r RedirectRender) Render(writer http.ResponseWriter, data interface{}) error {
	if r.Code == 0 {
		r.Code = http.StatusFound
	}
	http.Redirect(writer, r.Request, r.RedirectURL, r.Code)
	return nil
}
