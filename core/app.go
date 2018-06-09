package core

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Conf struct {
	AppName string `json:"appname"`
	Port    int    `json:"port"`
	AppRoot string `json:"approot"`
}

// RouteHandler callback func

type routerItem struct {
	method   string
	prefix   string
	handlers []gin.HandlerFunc
}

// App main
type App struct {
	Config       *Conf
	rootRouter   *gin.RouterGroup
	ginEngine    *gin.Engine
	routerGroups map[string]*gin.RouterGroup
	routers      map[string][]*routerItem
	midwares     map[string][]gin.HandlerFunc
}

// MidWare collect midwares
func (app *App) MidWare(
	group string,
	handler gin.HandlerFunc,
) {
	if _, ok := app.midwares[group]; !ok {
		app.midwares[group] = []gin.HandlerFunc{}
	}
	app.midwares[group] = append(app.midwares[group], handler)
}

// Init init vars
func (app *App) Init() {
	app.ginEngine = gin.Default()
	app.rootRouter = app.ginEngine.Group("/")
	app.routerGroups = map[string]*gin.RouterGroup{}
	app.routers = map[string][]*routerItem{}
	app.midwares = map[string][]gin.HandlerFunc{}
}

// Router collect and then make a router item
func (app *App) Router(
	group string,
	method string, // http verb list
	prefix string, // router
	handlers ...gin.HandlerFunc,
) {
	item := routerItem{
		method:   method,
		prefix:   prefix,
		handlers: handlers,
	}
	if _, ok := app.routers[group]; !ok {
		app.routers[group] = []*routerItem{}
	}
	app.routers[group] = append(app.routers[group], &item)

}

// SortedRouters asc sort by group length
func (app *App) SortedRouters() []string {
	_routers := [200]string{}
	for group := range app.routers { // base on routers
		_routers[len(group)] = group
	}
	routers := []string{}
	for _, router := range _routers {
		if router != "" {
			routers = append(routers, router)
		}
	}
	return routers
}

// AutoGroup make group logic match
func (app *App) AutoGroup(group string) *gin.RouterGroup {

	if routerGroup, ok := app.routerGroups[group]; ok {
		return routerGroup
	}

	groupStrs := strings.Split(group, "/")
	currentRouterGroup := app.rootRouter
	for _, currentGroupStr := range groupStrs {
		currentGroupStr = "/" + currentGroupStr
		if currentGroupStr == "/" {
			continue
		}
		if _, ok := app.routerGroups[currentGroupStr]; !ok {
			app.routerGroups[currentGroupStr] = currentRouterGroup.Group(currentGroupStr)
		}
		currentRouterGroup = app.routerGroups[currentGroupStr]
	}
	app.routerGroups[group] = currentRouterGroup
	return currentRouterGroup
}

// Start start app
func (app *App) Start() {
	sortedGroups := app.SortedRouters()
	for _, group := range sortedGroups {
		engine := app.AutoGroup(group)
		for _, midware := range app.midwares[group] {
			engine.Use(midware)
		}
		for _, router := range app.routers[group] {
			verbs := parseHTTPVerbs(router.method)
			ginHanlders := []gin.HandlerFunc{}
			// for _, handler := range router.handlers {
			// 	ginHanlders = append(ginHanlders, gin.HandlerFunc(handler))
			// }
			ginHanlders = router.handlers
			for _, method := range verbs {
				if method == "GET" {
					engine.GET(
						router.prefix,
						ginHanlders...,
					)
				} else if method == "POST" {
					engine.POST(
						router.prefix,
						ginHanlders...,
					)
				} else if method == "PUT" {
					engine.PUT(
						router.prefix,
						ginHanlders...,
					)
				} else if method == "DELETE" {
					engine.DELETE(
						router.prefix,
						ginHanlders...,
					)
				} else {
					panic(wrongMethodError{})
				}
			}
		}
	}

	if app.Config.Port == 0 {
		panic(noAddressError{})
	}

	app.ginEngine.Run(":" + strconv.Itoa(app.Config.Port))
}

func parseHTTPVerbs(method string) []string {
	methods := strings.Split(method, ",")
	return methods
}
