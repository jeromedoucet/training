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

	router.HandleFunc("/app/public/users", createUserHandlerFunc(c, conn))
	router.HandleFunc("/app/public/login", authenticationHandlerFunc(c, conn))

	router.ServeStaticAt("./front/dist/", route.Spa)

	return &trainingRouter{
		handler: router,
	}
}

type trainingRouter struct {
	handler http.Handler
}

func (r *trainingRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(res, req)
}
