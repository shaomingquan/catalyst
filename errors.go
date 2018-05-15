package core

type wrongMethodError struct{}

func (e *wrongMethodError) Error() string {
	return "wrong method, only GET,POST,PUT,DELETE httpverb supported"
}

type noAddressError struct{}

func (e *noAddressError) error() string {
	return "make sure address is setting in your app.Config"
}
