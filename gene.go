package main

import (
	"encoding/json"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	g "github.com/shaomingquan/webcore/gene"
)

var gene commandHandler

func init() {
	gene = func(command string, params []string) {
		routers, midwares, pkgs, ps, decorators, decoraorPkgs := getRoutersAndMidwares()

		pkgs = mergePkgs(pkgs, decoraorPkgs)

		// 1, update boot file
		updateBootFile()

		// 2, generate import file
		makeImportFile(routers, midwares, pkgs, ps, decorators)
	}
	handlers["gene"] = gene
}

var (
	pkgInfo *build.Package
)

func initBootFile() string {
	conf := getConf()
	rd := conf.AppRoot
	data := map[string]interface{}{
		"rootDir": rd,
	}
	content := g.MakeBootFile(data)
	return string(content)
}

func updateBootFile() {
	bootContentBt, err := ioutil.ReadFile(bootFileStr(rootRelatedString()))
	if err != nil {
		log.Fatal(err)
	}
	bootContent := string(bootContentBt)
	if bootContent == "" {
		bootContent = initBootFile()
	}
	pkgDir := getPkgDir()
	prefix := pkgDir[5:]
	if prefix == "" {
		prefix = "/"
	}
	appid := strings.Replace(prefix, "/", "_", -1)
	nextStr := `
		imports.Start` + appid + `(&app)
		// ###
	`
	bootContent = strings.Replace(bootContent, "// ###", nextStr, -1)
	ioutil.WriteFile(bootFileStr(rootRelatedString()), []byte(bootContent), 0644)
}

func mergePkgs(midwaresPkgs []map[string]string, decoratorsPkgs [][]map[string]string) []map[string]string {
	ret := midwaresPkgs
	checker := map[string]bool{}
	for _, pkg := range midwaresPkgs {
		checker[pkg["pkgid"]] = true
	}
	for _, pkgs := range decoratorsPkgs {
		for _, pkg := range pkgs {
			if !checker[pkg["pkgid"]] {
				ret = append(ret, pkg)
				checker[pkg["pkgid"]] = true
			}
		}
	}
	return ret
}

func makeImportFile(
	routers []string,
	midwares []map[string]string,
	pkgs []map[string]string,
	_params *map[string][][]string,
	_decorators *map[string][]map[string]string,
) {
	decorators := *_decorators
	params := *_params
	conf := getConf()
	rd := conf.AppRoot
	appName := pkgInfo.Name
	pkgDir := getPkgDir()
	prefix := pkgDir[5:]
	if prefix == "" {
		prefix = "/"
	}
	appid := strings.Replace(prefix, "/", "_", -1)
	filename := "app" + appid
	outputDir := importFileStr(rootRelatedString(), filename)

	routersWithValidator := []string{}
	routersNormal := []string{}

	for _, router := range routers {
		if _, ok := params[router]; ok {
			routersWithValidator = append(routersWithValidator, router)
		} else {
			routersNormal = append(routersNormal, router)
		}
	}

	hasvalidator := len(routersWithValidator) > 0

	data := map[string]interface{}{
		"routers":              routersNormal,
		"routersWithValidator": routersWithValidator,
		"hasvalidator":         hasvalidator,
		"appName":              appName,
		"pkgDir":               pkgDir, // subapi
		"prefix":               prefix,
		"appid":                appid,
		"midwares":             midwares,
		"decorators":           decorators,
		"pkgs":                 pkgs,
		"rootDir":              rd,
		"params":               params,
	}
	wirteMetaFile(data) // write meta file for debug GUI
	content := g.MakeImporterFile(data)
	ioutil.WriteFile(outputDir, content, 0644)
}

func wirteMetaFile(data map[string]interface{}) {
	filename := ".pkg_meta"
	datastring, err := json.Marshal(data)

	if err != nil {
		log.Fatal("json format wrong when trying to write meta data")
	}

	ioutil.WriteFile(filename, []byte(datastring), 0644)
}

func getRoutersAndMidwares() (
	[]string,
	[]map[string]string,
	[]map[string]string,
	*map[string][][]string,
	*map[string][]map[string]string,
	[][]map[string]string,
) {
	var err error
	pkgInfo, err = build.ImportDir(".", 0)
	if err != nil {
		log.Fatal(err)
	}
	fset := token.NewFileSet()
	gintance := g.Gene{}
	for _, file := range pkgInfo.GoFiles {
		if !g.IsUnderscoreFile(file) {
			continue
		}
		f, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			log.Fatal(err)
		}
		gintance.Collect(f)
	}
	a := gintance.OutputRouters()
	b, c := gintance.OutputMidwares()
	d := gintance.OutputParams()
	e, f := gintance.OutputDecorator()
	return a, b, c, d, e, f
}

var root string

func rootRelatedString() string {
	if root != "" {
		return root
	}
	_root := "../"
	level := 1
	for {
		_, err := ioutil.ReadFile(_root + "appconf.json") // appconf is required
		if level > 4 {
			log.Fatal("router too deep")
		}
		if err != nil {
			_root += "../"
			level++
			continue
		}
		break
	}
	root = _root
	return _root
}

func getPkgDir() string {
	absCurrent, _ := filepath.Abs(".")
	absRoot, _ := filepath.Abs(root)
	l := len(absRoot)
	if len(absCurrent) < l {
		log.Fatal("bad package dir")
	}
	ret := absCurrent[l:]
	return ret
}

func bootFileStr(root string) string {
	return root + "boot.go"
}

func importFileStr(root string, pkgname string) string {
	return root + "imports/" + pkgname + ".go"
}
