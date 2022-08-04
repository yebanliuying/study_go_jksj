package model

import (
	"time"

	"gorm.io/gorm"


)

//自定义公用字段
type BaseModel struct {
	ID int32 `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

type User struct {
	BaseModel
	Mobile string `gorm:"index:idx_mobile;type:varchar(11);unique;not null;comment:'手机号'"`
	Password string `gorm:"type:varchar(100);not null;comment:'密码'"`
	Nikename string `gorm:"varchar(20);comment:'昵称'"`
	Birthday *time.Time `gorm:"type:datetime;comment:'生日'"` //避免填充空时间
	Gender int `gorm:"type:tinyint(1);default:0;comment:'男:1,女:2,未知:0'"`
	Role int `gorm:"type:tinyint(2);default:1;comment:'普通用户:1,管理员:2'"`
}

func main() {

}
