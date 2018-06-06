package gene

import (
	"log"
	"strings"
)

type middwareItem struct {
	ID     string
	pkg    string
	method string
	params []string
}

type midwareItemSet []middwareItem

func (mset *midwareItemSet) collect(str string) {

	if !strings.Contains(str, "@") {
		log.Fatal("bad middware schema " + str)
	}

	coms := strings.Split(str, "@")
	if len(coms) < 2 {
		log.Fatal("bad middware schema " + str)
	}

	pkg := coms[0]

	if pkg == "" {
		log.Fatal("middware schema without pkgname " + str)
	}

	method := coms[1]

	if pkg == "" {
		log.Fatal("middware schema without method " + str)
	}

	_params := strings.Split(method, "#")
	params := []string{}

	if len(_params) == 2 {
		method = _params[0]
		params = strings.Split(_params[1], ",")
	}

	mid := strings.Replace(pkg, "/", "_", -1)
	if len(mid) > 0 && string(mid[0]) == "_" {
		mid = mid[1:]
	}

	*mset = append(*mset, middwareItem{mid, pkg, method, params})
}

func (mset *midwareItemSet) dump() ([]map[string]string, []map[string]string) {
	// package dedup
	pkgchecker := map[string]bool{}
	pkgList := []map[string]string{}
	for _, item := range *mset {
		if !pkgchecker[item.pkg] {
			pkgList = append(pkgList, map[string]string{
				"pkg":   item.pkg,
				"pkgid": item.ID,
			})
			pkgchecker[item.pkg] = true
		}
	}
	retset := []map[string]string{}
	for _, item := range *mset {
		params := make([]string, len(item.params))
		for index, p := range item.params {
			params[index] = "`" + p + "`" // why " occur error????
		}
		println(5555, strings.Join(params, ", "))
		retset = append(retset, map[string]string{
			"pkg":    item.pkg,
			"method": item.method,
			"pkgid":  item.ID,
			"params": strings.Join(params, ", "),
		})
	}
	return retset, pkgList
}
