package gene

import "strings"

type middwareItem struct {
	ID     string
	pkg    string
	method string
}

type midwareItemSet []middwareItem

func (mset *midwareItemSet) collect(str string) {
	coms := strings.Split(str, "@")
	pkg := coms[0]
	method := coms[1]
	mid := strings.Replace(pkg, "/", "_", -1)
	if len(mid) > 0 && string(mid[0]) == "_" {
		mid = mid[1:]
	}
	*mset = append(*mset, middwareItem{mid, pkg, method})
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
		retset = append(retset, map[string]string{
			"pkg":    item.pkg,
			"method": item.method,
			"pkgid":  item.ID,
		})
	}
	return retset, pkgList
}
