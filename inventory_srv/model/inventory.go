package model

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index"`
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"` // 分布式锁的乐观锁
}

// type InventoryHistory struct {
// 	User   int32
// 	Goods  int32
// 	Nums   int32
// 	Order  int32
// 	Status int32 // 1 表示库存是预扣减 2 表示已支付
// }
