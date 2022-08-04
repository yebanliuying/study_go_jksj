package global

import (
	ut "github.com/go-playground/universal-translator"
	"xshop_api/goods-web/config"
	"xshop_api/goods-web/proto"
)

var (
	Trans ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	GoodsSrvCln proto.GoodsClient
)