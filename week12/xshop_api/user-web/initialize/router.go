package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xshop_api/user-web/middlewares"
	"xshop_api/user-web/router"
)

func InitRouters() *gin.Engine {

	Router := gin.Default()

	Router.GET("/health", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"code":http.StatusOK,
			"success":true,
		})
	})

	//配置跨域
	Router.Use(middlewares.Cors())
	ApiGroup := Router.Group("/u/v1")
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup) //基本配置信息

	return Router
}