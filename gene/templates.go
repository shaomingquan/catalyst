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

import "github.com/gin-gonic/gin"
import core "github.com/shaomingquan/webcore/core"
{{if $.hasvalidator}}
import "net/http"
import validator "gopkg.in/validator.v2"
import gene "github.com/shaomingquan/webcore/gene"
{{end}}
import "{{.rootDir}}{{.pkgDir}}"
{{range $pkg := .pkgs}}
import {{$pkg.pkgid}} "{{$.rootDir}}/{{$pkg.pkg}}"
{{end}}

func Start{{.appid}}(app *core.App) {

	{{range $key, $item := .params}}
	var paramsValidatorOf{{$key}} = func(
		ctx *gin.Context,
	) {
		// validate
		paramsInstance := demo.ParamsOf{{$key}} {
			{{range $param := $item}}gene.ParamTo{{index $param 1}}(ctx, "{{index $param 0}}"),
			{{end}}
		}
		if err := validator.Validate(paramsInstance); err != nil {
			ctx.JSON(400, err.Error())
			ctx.AbortWithStatus(http.StatusBadRequest)
		} else {
			ctx.Next()
		}
	}
	{{end}}

	{{range $item := .midwares}}
	app.MidWare(
		"{{$.prefix}}", 
		{{$item.pkgid}}.{{$item.method}}({{$item.params}}),
	)
	{{end}}


	{{range $item := .routersWithValidator}}
	app.Router(
		"{{$.prefix}}",
		{{$.appName}}.MethodOf{{$item}}, 
		{{$.appName}}.PrefixOf{{$item}}, 
		{{$.appName}}.HandlerOf{{$item}},
		{{range $decorator := index $.decorators $item}}
		{{$decorator.pkgid}}.{{$decorator.method}}({{$decorator.params}}),
		{{end}}
		paramsValidatorOf{{$item}},
	)
	{{end}}

	{{range $item := .routers}}
	app.Router(
		"{{$.prefix}}",
		{{$.appName}}.MethodOf{{$item}}, 
		{{$.appName}}.PrefixOf{{$item}}, 
		{{$.appName}}.HandlerOf{{$item}},
		{{range $decorator := index $.decorators $item}}
		{{$decorator.pkgid}}.{{$decorator.method}}({{$decorator.params}}),
		{{end}}
		func(ctx *gin.Context) { ctx.Next() },
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
