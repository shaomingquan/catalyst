package gene

type routerItem struct {
	routerHandlerName bool
	routerPrefix      bool
	routerMethod      bool
}

type routerItemMap map[string]*routerItem

func (rm *routerItemMap) collect(name string, category string) {
	var item *routerItem
	if tmp, ok := (*rm)[name]; ok {
		item = tmp
	} else {
		item = &routerItem{}
		(*rm)[name] = item
	}
	if category == "prefix" {
		item.routerPrefix = true
	} else if category == "method" {
		item.routerMethod = true
	} else if category == "router" {
		item.routerHandlerName = true
	}
}

func (rm *routerItemMap) dump(fn *func(name string)) {
	rrm := *rm
	for name := range rrm {
		m := rrm[name]
		if m.routerHandlerName &&
			m.routerPrefix &&
			m.routerMethod {
			(*fn)(name)
		} else {
			panic(routerItemIncompleteError{"", name})
		}
	}
}
