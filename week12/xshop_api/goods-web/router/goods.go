package router

import (
	"github.com/gin-gonic/gin"
	"xshop_api/goods-web/api/goods"
	"xshop_api/goods-web/middlewares"
)

func InitGoodsRouter(Router *gin.RouterGroup)  {
	//UserRouter := Router.Group("user").Use(middlewares.JWTAuth())
	//加上链路追踪
	GoodsRouter := Router.Group("goods").Use(middlewares.Trace())
	{
		GoodsRouter.GET("", goods.List)
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)
		GoodsRouter.GET("/:id", goods.Detail)
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Detail)
		GoodsRouter.GET("/:id/stocks", goods.Stocks) //获取商品库存

		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus)

	}
}