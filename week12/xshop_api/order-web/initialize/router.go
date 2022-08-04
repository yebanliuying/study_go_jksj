package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xshop_api/order-web/middlewares"
	"xshop_api/order-web/router"
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
	ApiGroup := Router.Group("/o/v1")
	router.InitGoodsRouter(ApiGroup)
	//router.InitCategoryRouter(ApiGroup)
	//router.InitBannerRouter(ApiGroup)
	//router.InitBrandRouter(ApiGroup)

	return Router
}