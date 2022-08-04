package global

import (
	ut "github.com/go-playground/universal-translator"
	"xshop_api/order-web/config"
	"xshop_api/order-web/proto"
)

var (
	Trans ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	GoodsSrvCln proto.GoodsClient

	OrderSrvCln proto.OrderClient

	InventorySrvCln proto.InventoryClient
)