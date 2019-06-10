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
	dRouter := route.NewDynamicRouter()

	dRouter.HandleFunc("/app/users", createUserHandlerFunc(c, conn))
	dRouter.HandleFunc("/app/login", authenticationHandlerFunc(c, conn))

	return &trainingRouter{
		xhrHandler:    dRouter,
		staticHandler: http.FileServer(http.Dir("front/dist")),
	}
}

type trainingRouter struct {
	xhrHandler    *route.DynamicRouter
	staticHandler http.Handler
}

func (r *trainingRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if strings.Contains(req.URL.Path, "app/") {
		r.xhrHandler.ServeHTTP(res, req)
	} else {
		r.staticHandler.ServeHTTP(res, req)
	}

}
