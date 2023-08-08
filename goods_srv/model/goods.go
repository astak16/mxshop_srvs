package model

type Category struct {
	BaseModel
	Name             string      `gorm:"type:varchar(20);not null" json:"name"`
	ParentCategoryId int32       `json:"parent"`
	ParentCategory   *Category   `json:"-"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryId;references:ID" json:"sub_category"`
	Level            int32       `gorm:"type:int;not null;default:1" json:"level"`
	IsTab            bool        `gorm:"not null;default:false" json:"is_tab"`
}

type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(100);not null;default:''"`
}

type GoodsCategoryBrand struct {
	BaseModel
	CategoryId int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category
	BrandsId   int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brands     Brands
}

func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index int32  `gorm:"type:int;not null;default:1"`
}

type Goods struct {
	BaseModel

	CategoryId int32 `gorm:"type:int;not null"`
	Category   Category
	BrandsId   int32 `gorm:"type:int;not null"`
	Brands     Brands

	OnSale   bool `gorm:"not null;default:false"`
	ShipFree bool `gorm:"not null;default:false"`
	IsNew    bool `gorm:"not null;default:false"`
	IsHot    bool `gorm:"not null;default:false"`

	Name            string   `gorm:"type:varchar(50);not null"`
	GoodsSn         string   `gorm:"type:varchar(50);not null"`
	ClickNum        int32    `gorm:"type:int;not null;default:0"`
	SoldNum         int32    `gorm:"type:int;not null;default:0"`
	FavNum          int32    `gorm:"type:int;not null;default:0"`
	MarketPrice     float32  `gorm:"not null"`
	ShopPrice       float32  `gorm:"not null"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null"`
	Images          GormList `gorm:"type:varchar(1000);not null"`
	DescImages      GormList `gorm:"type:varchar(1000);not null"`
	GoodsFrontImage string   `gorm:"type:varchar(100);not null"`
}
