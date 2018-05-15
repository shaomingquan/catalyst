package gene

type routerItemIncompleteError struct {
	which string
	param string
}

func (e *routerItemIncompleteError) Error() string {
	return "in routeItem " + e.which + " param " + e.param + " required"
}
