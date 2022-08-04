package global

import (
	"log"
	"os"
	"time"
	"xshop_srvs/user_srv/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	
)

var (
	DB *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig config.NacosConfig
)

func init()  {
	dsn := "root:root@tcp(139.198.21.42:3306)/xshop_user_server?charset=utf8mb4&parseTime=True&loc=Local"


	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold: time.Second,   // 慢 SQL 阈值
			LogLevel:      logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,   // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:      true,         // 彩色打印
		},
	)

	// 全局模式
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "xs_",//表前缀
			SingularTable: true,//不走默认复数表名
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

}