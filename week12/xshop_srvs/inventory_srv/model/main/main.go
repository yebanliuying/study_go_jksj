package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
	"xshop_srvs/inventory_srv/model"
)

func main() {
	dsn := "root:root@tcp(139.198.21.42:3306)/xshop_inventory_server?charset=utf8mb4&parseTime=True&loc=Local"


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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "xs_",//表前缀
			SingularTable: true,//不走默认复数表名
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}



	//定义表结构
	//_ = db.AutoMigrate(&model.Inventory{}, &model.StockSellDetail{}) //生成表

	//orderDetail := model.StockSellDetail{
	//	OrderSn: "laogeceshi",
	//	Status: 1,
	//	Detail: []model.GoodsDetail{{1,2},{2,3}},
	//}
	//
	//db.Create(&orderDetail)
	
	var sellDetail model.StockSellDetail
	db.Where(model.StockSellDetail{OrderSn:"laogeceshi"}).First(&sellDetail)
	fmt.Println(sellDetail.Detail)

}