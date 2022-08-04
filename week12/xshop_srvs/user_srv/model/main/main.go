package main

import (
	"crypto/sha512"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
	"xshop_srvs/user_srv/model"
	"xshop_srvs/user_srv/utils"
)

func main() {
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

	//密码加密
	options := &utils.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := utils.Encode("123456", options)
	userPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s",salt,encodedPwd)

	for i := 0; i < 10; i++ {
		user := model.User{
			Nikename: fmt.Sprintf("老弟_%d", i),
			Mobile: fmt.Sprintf("135111111%d", i),
			Password: userPassword,
		}
		db.Save(&user)
	}


	//定义表结构
	//迁移 schema
	//_ = db.AutoMigrate(&model.User{}) //生成表




	//// Create
	//db.Create(&Product{Code: "D42", Price: 100})
	//
	//// Read
	//var product Product
	//db.First(&product, 1) // 根据整型主键查找
	//db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
	//
	//// Update - 将 product 的 price 更新为 200
	//db.Model(&product).Update("Price", 200)
	//// Update - 更新多个字段
	//db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	////db.Model(&product).Updates(Product{Price: 200, Code: sql.NullString{}}) // 仅更新非零值字段
	//db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
	//
	//// Delete - 删除 product
	//db.Delete(&product, 1)

}