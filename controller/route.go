package controller

import (
	"net/http"
	"strings"

	"github.com/jeromedoucet/route"
	"github.com/jeromedoucet/training/configuration"
	"github.com/jeromedoucet/training/dao"
)

func InitRoutes(c *configuration.GlobalConf) http.Handler {
	conn := dao.Open(c)
	pRouter := route.NewDynamicRouter()

	pRouter.HandleFunc("/app/public/users", createUserHandlerFunc(c, conn))
	pRouter.HandleFunc("/app/public/login", authenticationHandlerFunc(c, conn))

	return &trainingRouter{
		xhrHandler:    pRouter,
		staticHandler: http.FileServer(http.Dir("front/dist")),
	}
}

type trainingRouter struct {
	xhrHandler    http.Handler
	staticHandler http.Handler
}

func (r *trainingRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if strings.Contains(req.URL.Path, "app/") {
		r.xhrHandler.ServeHTTP(res, req)
	} else {
		r.staticHandler.ServeHTTP(res, req)
	}
}
