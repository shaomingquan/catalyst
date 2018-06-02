package gene

import (
	"go/ast"
)

type traverseFunc func(string, string, interface{})

// traverse width only first level, and fast reture
func traverse(current ast.Node, fn traverseFunc) ast.Node {
	switch n := current.(type) {
	case *ast.File:
		walkDeclList(n.Decls, fn)
	case *ast.GenDecl:
		// 取出路由定义和方法定义
		name, appends := parseNameFromGenDecl(n)
		if name != "" {
			fn(name, "var", appends)
		}
	case *ast.FuncDecl:
		// 取出路由
		name := parseNameFromFuncDecl(n)
		if name != "" {
			fn(name, "fn", "")
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
func parseNameFromGenDecl(node *ast.GenDecl) (string, interface{}) {
	specVal, ok := node.Specs[0].(*ast.ValueSpec)
	if !ok {
		return "", nil
	}

	name := specVal.Names[0].Name
	if name == "MiddlewaresComposer" { // middware decl
		list := specVal.Values[0].(*ast.CompositeLit)
		midwares := []string{}
		for _, exp := range list.Elts {
			importStr := exp.(*ast.BasicLit).Value
			if importStr == "" {
				continue
			}
			l := len(importStr)
			importStr = importStr[:l-1]
			importStr = importStr[1:]
			midwares = append(midwares, importStr)
		}
		return name, midwares
	} else {
		return name, nil
	}
}

func parseNameFromFuncDecl(node *ast.FuncDecl) string {
	return node.Name.Name
}
