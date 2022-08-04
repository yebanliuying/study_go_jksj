package global

import (
	ut "github.com/go-playground/universal-translator"
	"xshop_api/user-web/config"
	"xshop_api/user-web/proto"
)

var (
	Trans ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	UserSrvCln proto.UserClient
)