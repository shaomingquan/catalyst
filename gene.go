package main

import (
	"flag"
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
		routers, midwares, pkgs := getRoutersAndMidwares()

		// 1, update boot file
		updateBootFile()

		// 2, generate import file
		makeImportFile(routers, midwares, pkgs)
	}
	handlers["gene"] = gene
}

var (
	pkgInfo *build.Package
	stage   = flag.String("stage", "main", "run stage")
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

func makeImportFile(
	routers []string,
	midwares []map[string]string,
	pkgs []map[string]string,
) {
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
	data := map[string]interface{}{
		"routers":  routers,
		"appName":  appName,
		"pkgDir":   pkgDir, // subapi
		"prefix":   prefix,
		"appid":    appid,
		"midwares": midwares,
		"pkgs":     pkgs,
		"rootDir":  rd,
	}
	content := g.MakeImporterFile(data)
	ioutil.WriteFile(outputDir, content, 0644)
}

func getRoutersAndMidwares() ([]string, []map[string]string, []map[string]string) {
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
	return a, b, c
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
