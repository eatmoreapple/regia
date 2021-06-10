package regia

import "net/http"

type Exit interface{ Exit(context *Context) }

type exit struct{}

// Exit Do nothing
func (e exit) Exit(*Context) {}

type AuthenticationFailed struct{}

func (a AuthenticationFailed) Exit(context *Context) {
	http.Error(context.ResponseWriter, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func (a AuthenticationFailed) Error() string {
	return http.StatusText(http.StatusUnauthorized)
}

type ParseError struct{}

func (p ParseError) Exit(context *Context) {
	http.Error(context.ResponseWriter, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func (p ParseError) Error() string {
	return http.StatusText(http.StatusBadRequest)
}

type PermissionDenied struct{}

func (p PermissionDenied) Exit(context *Context) {
	http.Error(context.ResponseWriter, http.StatusText(http.StatusForbidden), http.StatusForbidden)
}

func (p PermissionDenied) Error() string {
	return http.StatusText(http.StatusForbidden)
}
