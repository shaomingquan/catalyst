package gene

import (
	"go/ast"
	"go/token"
	"strings"
)

// Gene verbose
type Gene struct {
	result []string
}

var (
	fset     = token.NewFileSet()
	routers  = routerItemMap{}
	midwares = midwareItemSet{}
	params   = paramsMap{}
)

// Collect every _ file do this
func (g *Gene) Collect(f *ast.File) {
	traverseCallback := func(name, t string, appends interface{}) {

		if len(name) <= 8 { // seems all less than 8
			return
		}

		if name == "MiddlewaresComposer" {
			mids := appends.([]string)
			for _, mid := range mids {
				midwares.collect(mid)
			}
			return
		}

		if name[:8] == "ParamsOf" {
			params.collect(name[8:], appends.([][]string))
		}

		if strings.HasPrefix(name, "PrefixOf") {
			routers.collect(name[8:], "prefix")
		} else if strings.HasPrefix(name, "MethodOf") {
			routers.collect(name[8:], "method")
		} else if strings.HasPrefix(name, "HandlerOf") {
			routers.collect(name[9:], "router")
		}
	}
	traverse(f, traverseCallback)
}

// OutputRouters _ files output the routers togother
func (g *Gene) OutputRouters() []string {
	rs := []string{}
	dumpCallback := func(name string) {
		rs = append(rs, name)
	}
	routers.dump(&dumpCallback)
	return rs
}

// OutputMidwares output midware ast result
func (g *Gene) OutputMidwares() ([]map[string]string, []map[string]string) {
	return midwares.dump()
}

// OutputParams params of each router map[router]params[name, type]
func (g *Gene) OutputParams() *map[string][][]string {
	return params.dump()
}

// IsUnderscoreFile is _ file
func IsUnderscoreFile(filename string) bool {
	return strings.HasPrefix(filename, "t.")
}
