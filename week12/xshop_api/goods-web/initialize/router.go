package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xshop_api/goods-web/middlewares"
	"xshop_api/goods-web/router"
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
	ApiGroup := Router.Group("/g/v1")
	//ApiGroup := Router.Group("/v1")
	router.InitGoodsRouter(ApiGroup)
	router.InitCategoryRouter(ApiGroup)
	router.InitBannerRouter(ApiGroup)
	router.InitBrandRouter(ApiGroup)

	return Router
}