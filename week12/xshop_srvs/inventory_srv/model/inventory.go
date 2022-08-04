package model

import (
	"database/sql/driver"
	"encoding/json"
)

//type Stock struct {
//	BaseModel
//	Name    string
//	Address string
//}

type GoodsDetail struct {
	Goods int32
	Num int32
}
type GoodsDetailList []GoodsDetail

//实现 driver.Valuer 接口，Value 返回 json value
func (g GoodsDetailList) Value() (driver.Value, error){
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GoodsDetailList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

type Inventory struct {
	BaseModel
	Goods  int32 `gorm:"type:int;index;"`
	Stocks int32 `gorm:"type:int"`
	//Stock Stock
	Version int32 `gorm:"type:int"` //分布式锁的乐观锁

}

type Delivery struct {
	Goods int32 `gorm:"type:int;index;"`
	Nums int32 `gorm:"type:int"`
	OrderSn string `gorm:"type:varchar(200)"`
	Status int32 `gorm:"type:int"` //1、已扣减 2、已归还
}

type StockSellDetail struct{
	OrderSn string `gorm:"type:varchar(200);index:idx_order_sn,unique;"`
	Status int32 `gorm:"type:int"`  //1、已扣减 2、已归还
	Detail GoodsDetailList `gorm:"type:varchar(200)"`
}

//type InventoryHistory struct {
//	user int32
//	goods int32
//	nums int32
//	order int32
//	status int32 //1、库存预扣减 幂等性 2、已经支付成功
//}
