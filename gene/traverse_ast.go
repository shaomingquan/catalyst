package gene

import "go/ast"

type traverseFunc func(string, string)

// traverse width only first level, and fast reture
func traverse(current ast.Node, fn traverseFunc) ast.Node {
	switch n := current.(type) {
	case *ast.File:
		walkDeclList(n.Decls, fn)
	case *ast.GenDecl:
		// 取出路由定义和方法定义
		name := parseNameFromGenDecl(n)
		if name != "" {
			fn(name, "var")
		}
	case *ast.FuncDecl:
		// 取出路由
		name := parseNameFromFuncDecl(n)
		if name != "" {
			fn(name, "fn")
		}
	}
	return current
}

func walkDeclList(list []ast.Decl, fn traverseFunc) {
	for _, x := range list {
		traverse(x, fn)
	}
}

// 暂时无法识别多赋值
func parseNameFromGenDecl(node *ast.GenDecl) string {
	specVal, ok := node.Specs[0].(*ast.ValueSpec)
	if !ok {
		return ""
	}
	kind := specVal.Values[0].(*ast.BasicLit).Kind
	if kind.String() != "STRING" {
		return ""
	}
	return specVal.Names[0].Name

}

func parseNameFromFuncDecl(node *ast.FuncDecl) string {
	return node.Name.Name
}
