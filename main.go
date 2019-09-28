package main

import (
	"github.com/susengo/commontools/xorm"
	"github.com/susengo/walletdc/config"
	"github.com/susengo/walletdc/controller"
	"github.com/susengo/walletdc/service/handler"
	"github.com/susengo/walletdc/service/router"
)

func main() {

	configFile := config.Config.GetString("SwingGatewayDbconfig")
	port := config.Config.GetString("SwingGatewayPort")

	dbengine := xorm.GetEngine(configFile)
	apiController := controller.NewApiController(
		handler.NewUserHandler(dbengine),
		handler.NewRoleHandler(dbengine),
	)

	router.SetRouterAndRun(apiController, port)

}
