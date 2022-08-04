package model

type Category struct {
	BaseModel
	Name             string `gorm:"type:varchar(20);not null;comment:'分类名'" json:"name"`
	Level            int32  `gorm:"type:int;not null;default:1;comment:'分类等级'" json:"level"`
	ParentCategoryID int32 `json:"parent"`
	ParentCategory   *Category `json:"-"`
	SubCategory []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
	IsTab            bool `gorm:"default:false;not null" json:"is_tab"`
}

type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null;comment:'品牌名称'"`
	Logo string `gorm:"type:varchar(200);not null;default:'';comment:'品牌logo'"`
}

type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique;"`
	Category   Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand,unique;"`
	Brands   Brands
}

//重载定义分配品牌关联表名
//func (GoodsCategoryBrand) TableName() string {
//	return "goodscategorybrand"
//}

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null;default:'';comment:'图片地址';"`
	Url   string `gorm:"type:varchar(200);not null;default:'';comment:'链接地址';"`
	Index int32  `gorm:"type:int;not null;default:100;comment:'排序';"`
}

type Goods struct {
	BaseModel

	CategoryID int32 `gorm:"type:int;not null;"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;not null;"`
	Brands     Brands

	OnSale   bool `gorm:"default:false;not null;"`
	ShipFree bool `gorm:"default:false;not null;"`
	IsNew    bool `gorm:"default:false;not null;"`
	IsHot    bool `gorm:"default:false;not null;"`

	Name            string   `gorm:"type:varchar(50);not null;"`
	GoodsSn         string   `gorm:"type:varchar(50);not null;comment:'商品编码';"`
	ClickNum        int32    `gorm:"type:int;not null;default:0;comment:'点击数';"`
	SoldNum         int32    `gorm:"type:int;default:0;not null;comment:'已销售数量';"`
	FavNum          int32    `gorm:"type:int;default:0;not null;comment:'收藏数量';"`
	MarketPrice     float32  `gorm:"default:0;not null;comment:'市场价（划线价）';"`
	ShopPrice       float32  `gorm:"default:0;not null;comment:'售卖价格';"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null;comment:'商品简介';"`
	Images          GormList `gorm:"type:varchar(1000);not null;comment:'商品头图';"`
	DescImages      GormList `gorm:"type:varchar(1000);not null;comment:'商品详情图';"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null;comment:'图片封面图';"`
}
