package gene

type decoratorMap map[string]midwareItemSet

func (dmap *decoratorMap) collect(router string, str string) {
	if _, ok := (*dmap)[router]; !ok {
		(*dmap)[router] = midwareItemSet{}
	}
	midItem := parseMiddware(str)
	(*dmap)[router] = append((*dmap)[router], *midItem)
}

func (dmap *decoratorMap) dump() (map[string][]map[string]string, [][]map[string]string) {
	retset := map[string][]map[string]string{}
	pkglist := [][]map[string]string{}
	for router, mset := range *dmap {
		if _, ok := retset[router]; !ok {
			retset[router] = []map[string]string{}
		}
		retsetForRouter, pkgListForRouter := mset.dump()
		retset[router] = retsetForRouter
		pkglist = append(pkglist, pkgListForRouter)
	}
	return retset, pkglist
}
