package global

import (
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"

	"xshop_srvs/order_srv/config"
	"xshop_srvs/order_srv/proto"
)

var (
	DB           *gorm.DB
	RS           *redsync.Redsync
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig

	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
)

func init() {
	//dsn := "root:root@tcp(139.198.21.42:3306)/xshop_user_server?charset=utf8mb4&parseTime=True&loc=Local"
	//
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
	//	logger.Config{
	//		SlowThreshold:             time.Second, // 慢 SQL 阈值
	//		LogLevel:                  logger.Info, // 日志级别
	//		IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
	//		Colorful:                  true,        // 彩色打印
	//	},
	//)

	//// mysql
	//var err error
	//DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	NamingStrategy: schema.NamingStrategy{
	//		TablePrefix:   "xs_", //表前缀
	//		SingularTable: true,  //不走默认复数表名
	//	},
	//	Logger: newLogger,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(123123123)
	//
	//fmt.Println(ServerConfig.RedisInfo.Host)
	//fmt.Println(ServerConfig.RedisInfo.Port)
	//
	////redsync redis分布式锁
	//client := goredislib.NewClient(&goredislib.Options{
	//	Addr: fmt.Sprintf("%s:%d", ServerConfig.RedisInfo.Host, ServerConfig.RedisInfo.Port),
	//})
	//
	//pool := goredis.NewPool(client)
	//
	//RS = redsync.New(pool)
}
