package main

import (
	"fmt"
)

var help commandHandler

func init() {
	help = func(command string, params []string) {

		fmt.Println(`
hi!


[new]    =>    create new project
	example: webcore new your-project-name -port=8080 -tpl=curd
		-port: app default start port
		-tpl: app start template
			-tpl=curd: curd base tpl


[dev]    =>    run dev server
	example: webcore dev
		-withoutrun: only gene bootfile
		-port: start port


[update] =>    update webcore runtime lib and clitool
	example: webcore update
		`)

	}
	handlers["help"] = help
}
