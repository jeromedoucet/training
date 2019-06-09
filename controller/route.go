package controller

import (
	"net/http"

	"github.com/jeromedoucet/route"
	"github.com/jeromedoucet/training/configuration"
	"github.com/jeromedoucet/training/dao"
)

func InitRoutes(c *configuration.GlobalConf) http.Handler {
	conn := dao.Open(c)
	router := route.NewDynamicRouter()
	router.HandleFunc("/app/users", createUserHandlerFunc(c, conn))
	router.HandleFunc("/app/login", authenticationHandlerFunc(c, conn))
	return router
}
