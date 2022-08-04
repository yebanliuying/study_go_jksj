package router

import (
	"github.com/gin-gonic/gin"
	"xshop_api/goods-web/api/goods"
	"xshop_api/goods-web/middlewares"
)

func InitCategoryRouter(Router *gin.RouterGroup)  {
	//UserRouter := Router.Group("user").Use(middlewares.JWTAuth())
	GoodsRouter := Router.Group("categorys").Use(middlewares.Trace())
	{
		GoodsRouter.GET("", goods.List)
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)
		GoodsRouter.GET("/:id", goods.Detail)
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Detail)
		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)

	}
}