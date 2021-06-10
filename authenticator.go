package regia

type Authenticator interface {
	Authenticate(context *Context, v interface{}) (ok bool, err error)
}

type Authenticators []Authenticator

func (a Authenticators) RunAuthenticate(context *Context, v interface{}) error {
	for _, authenticator := range a {
		ok, err := authenticator.Authenticate(context, v)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}
	return AuthenticationFailed{}
}
