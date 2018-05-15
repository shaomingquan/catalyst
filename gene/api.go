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
	fset       = token.NewFileSet()
	collection = routerItemMap{}
)

// Collect every _ file do this
func (g *Gene) Collect(f *ast.File) {
	traverseCallback := func(name, t string) {
		if strings.HasPrefix(name, "PrefixOf") {
			collection.collect(name[8:], "prefix")
		} else if strings.HasPrefix(name, "MethodOf") {
			collection.collect(name[8:], "method")
		} else if strings.HasPrefix(name, "HandlerOf") {
			collection.collect(name[9:], "router")
		}
	}
	traverse(f, traverseCallback)
}

// OutputRouters _ files output the routers togother
func (g *Gene) OutputRouters() []string {
	routers := []string{}
	dumpCallback := func(name string) {
		routers = append(routers, name)
	}
	collection.dump(&dumpCallback)
	return routers
}

// IsUnderscoreFile is _ file
func IsUnderscoreFile(filename string) bool {
	return strings.HasPrefix(filename, "t.")
}
