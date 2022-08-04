package router

import (
	"github.com/gin-gonic/gin"
	"xshop_api/user-web/api"
)

//InitUserRouter 用户路由
func InitUserRouter(Router *gin.RouterGroup)  {
	//UserRouter := Router.Group("user").Use(middlewares.JWTAuth())
	UserRouter := Router.Group("user")
	{
		UserRouter.GET("list", api.GetUserList)
		//UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("pwd_login", api.PasswordLogin)
		UserRouter.POST("register", api.Register)
	}
}