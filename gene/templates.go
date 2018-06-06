package gene

import (
	"bytes"
	"go/format"
	"html/template"
	"log"
	"strings"
)

var bootTpl = `package main

import "{{.rootDir}}/imports"

func init() {
    go func() {
        //la prepare // wait for prepare
        app.Init()
    
        // ###

        loaded //la 1 // i am loaded
    }()

}

// auto generate by _, dont modify`

var importerTpl = `package imports

import core "github.com/shaomingquan/webcore/core"
import "{{.rootDir}}{{.pkgDir}}"
{{range $pkg := .pkgs}}
import {{$pkg.pkgid}} "{{$.rootDir}}/{{$pkg.pkg}}"
{{end}}

func Start{{.appid}}(app *core.App) {
	{{range $item := .midwares}}
	app.MidWare("{{$.prefix}}", {{$item.pkgid}}.{{$item.method}}({{$item.params}}))
	{{end}}

	{{range $item := .routers}}
	app.Router(
		"{{$.prefix}}",
		{{$.appName}}.MethodOf{{$item}}, 
		{{$.appName}}.PrefixOf{{$item}}, 
		{{$.appName}}.HandlerOf{{$item}},
	)
	{{end}}
}

// auto generate by _, dont modify`

func tplCommon(data map[string]interface{}, tpl string) []byte {
	t, err := template.New("").Parse(tpl)
	if err != nil {
		log.Fatal(err)
	}
	buff := bytes.NewBufferString("")
	err = t.Execute(buff, data)
	if err != nil {
		log.Fatal(err)
	}
	ret, err := format.Source(buff.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	return ret
}

// MakeBootFile boot file template
func MakeBootFile(data map[string]interface{}) []byte {
	f := tplCommon(data, bootTpl)
	filestring := string(f)
	filestring = strings.Replace(filestring, "//la", "<-", -1)
	return []byte(filestring)
}

// MakeImporterFile importer file template
func MakeImporterFile(data map[string]interface{}) []byte {
	return tplCommon(data, importerTpl)
}
