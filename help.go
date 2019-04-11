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
	example: catalyst new your-project-name -port=8080 -tpl=crud
		-port: app default start port
		-tpl: app start template
			-tpl=crud: crud base tpl


[dev]    =>    run dev server
	example: catalyst dev
		-withoutrun: only gene bootfile
		-port: start port


[update] =>    update catalyst runtime lib and clitool
	example: catalyst update
		`)

	}
	handlers["help"] = help
}
