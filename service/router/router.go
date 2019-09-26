package router

import (
	"github.com/gin-gonic/gin"
	"github.com/susengo/commontools/gintool"
	"github.com/susengo/swing/controller"
)

func SetRouterAndRun(api_controller *controller.ApiController, port string) {

	r := gin.New()
	r.Use(gintool.Logger())
	r.Use(gin.Recovery())

	gintool.UseSession(r)

	api := r.Group("/api")
	{

		api.POST("/user/login", api_controller.UserLogin)
		api.POST("/user/logout", api_controller.UserLogout)
		//认证校验
		api.Use(api_controller.UserAuthorize)
		api.GET("/user/info", api_controller.UserInfo)
		api.GET("/user/list", api_controller.UserList)
		api.POST("/user/add", api_controller.UserAdd)
		api.POST("/user/addAuth", api_controller.UserAddAuth)
		api.POST("/user/delAuth", api_controller.UserDelAuth)
		api.POST("/user/update", api_controller.UserUpdate)
		api.POST("/user/delete", api_controller.UserDelete)

	}

	r.Run(":" + port)
}
