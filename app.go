package core

import (
	"github.com/gin-gonic/gin"
)

// RouteHandler callback func
type RouteHandler func(ctx *gin.Context)

type routerItem struct {
	method  string
	prefix  string
	handler RouteHandler
}

// App main
type App struct {
	Name      string
	Config    map[string]string
	ginEngine *gin.Engine
	routers   []*routerItem
}

// Collect collect and then make a router item
func (app *App) Collect(
	method string, // http verb
	prefix string, // router
	handler RouteHandler, // handler
) {
	item := routerItem{
		method:  method,
		prefix:  prefix,
		handler: handler,
	}
	if app.ginEngine == nil {
		app.ginEngine = gin.Default()
	}
	if app.routers == nil {
		app.routers = []*routerItem{}
	}
	app.routers = append(app.routers, &item)

}

// Start start app
func (app *App) Start() {
	engine := app.ginEngine
	// TODO: 安装路由
	for _, router := range app.routers {
		if router.method == "GET" {
			engine.GET(
				router.prefix,
				gin.HandlerFunc(router.handler),
			)
		} else if router.method == "POST" {
			engine.POST(
				router.prefix,
				gin.HandlerFunc(router.handler),
			)
		} else if router.method == "PUT" {
			engine.PUT(
				router.prefix,
				gin.HandlerFunc(router.handler),
			)
		} else if router.method == "DELETE" {
			engine.DELETE(
				router.prefix,
				gin.HandlerFunc(router.handler),
			)
		} else {
			panic(wrongMethodError{})
		}
	}

	if app.Config["address"] == "" {
		panic(noAddressError{})
	}

	engine.Run(app.Config["address"])
}
