package model

import "time"


type ShoppingCart struct{
	BaseModel
	User int32 `gorm:"type:int;index"` //在购物车列表中需要查询当前用户的购物车记录
	Goods int32 `gorm:"type:int;index"`
	Nums int32 `gorm:"type:int"`
	Checked bool //是否选中
}

//func (ShoppingCart) TableName() string {
//	return "shoppingcart"
//}

type OrderInfo struct{
	BaseModel

	User int32 `gorm:"type:int;index"`
	OrderSn string `gorm:"type:varchar(30);index"` //订单号
	PayType string `gorm:"type:varchar(20) comment 'alipay(支付宝)， wechat(微信)'"`

	//status也可以使用iota
	Status string `gorm:"type:varchar(20)  comment 'PAYING(待支付), TRADE_SUCCESS(成功)， TRADE_CLOSED(超时关闭), WAIT_BUYER_PAY(交易创建), TRADE_FINISHED(交易结束)'"`
	TradeNo string `gorm:"type:varchar(100) comment '交易号'"` //第三方交易号
	OrderMount float32
	PayTime *time.Time `gorm:"type:datetime"`

	Address string `gorm:"type:varchar(100)"`
	SignerName string `gorm:"type:varchar(20)"`
	SingerMobile string `gorm:"type:varchar(11)"`
	Post string `gorm:"type:varchar(20)"`
}

//func (OrderInfo) TableName() string {
//	return "orderinfo"
//}

type OrderGoods struct{
	BaseModel

	Order int32 `gorm:"type:int;index"`
	Goods int32 `gorm:"type:int;index"`

	//高并发系统中一般都不会遵循三范式  做镜像记录
	GoodsName string `gorm:"type:varchar(100);index"`
	GoodsImage string `gorm:"type:varchar(200)"`
	GoodsPrice float32
	Nums int32 `gorm:"type:int"`
}

//func (OrderGoods) TableName() string {
//	return "ordergoods"
//}