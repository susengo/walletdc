package main

import (
	"github.com/susengo/commontools/xorm"
	"github.com/susengo/swing/config"
	"github.com/susengo/swing/controller"
	"github.com/susengo/swing/service/router"
	"github.com/susengo/swing/service/server"
)

func main() {

	configFile := config.Config.GetString("SwingGatewayDbconfig")
	port := config.Config.GetString("SwingGatewayPort")

	dbengine := xorm.GetEngine(configFile)
	apiController := controller.NewApiController(
		server.NewUserService(dbengine),
		server.NewRoleService(dbengine),
	)

	router.SetRouterAndRun(apiController, port)

}
