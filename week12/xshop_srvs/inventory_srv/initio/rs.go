package initio

import (
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"xshop_srvs/inventory_srv/global"
)

func InitRS() {

	//redsync redis分布式锁
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})

	pool := goredis.NewPool(client)

	global.RS = redsync.New(pool)

}