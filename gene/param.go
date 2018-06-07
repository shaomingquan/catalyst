package gene

type paramsMap map[string][][]string

func (p *paramsMap) collect(router string, params [][]string) {
	if _, ok := (*p)[router]; !ok {
		(*p)[router] = params
	}
}

func (p *paramsMap) dump() *map[string][][]string {
	ret := map[string][][]string(*p)
	return &ret
}
