package router

import (
	"github.com/gin-gonic/gin"
	"github.com/susengo/commontools/gintool"
	"github.com/susengo/walletdc/controller"
)

func SetRouterAndRun(ac *controller.ApiController, port string) {

	r := gin.New()
	r.Use(gintool.Logger())
	r.Use(gin.Recovery())

	gintool.UseSession(r)

	api := r.Group("/api")
	{

		api.POST("/user/login", ac.UserLogin)
		api.POST("/user/logout", ac.UserLogout)
		//认证校验
		api.Use(ac.UserAuthorize)
		api.GET("/user/info", ac.UserInfo)
		api.GET("/user/list", ac.UserList)
		api.POST("/user/add", ac.UserAdd)
		api.POST("/user/addAuth", ac.UserAddAuth)
		api.POST("/user/delAuth", ac.UserDelAuth)
		api.POST("/user/update", ac.UserUpdate)
		api.POST("/user/delete", ac.UserDelete)

	}

	r.Run(":" + port)
}
