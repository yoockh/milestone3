package entity

type Payment struct {
	Id int
	UserId int
	User Users `gorm:"foreignKey:UserId;references:Id"`
	AuctionItemId int
	Status string
	OrderId string
	// PaymentStatus PaymentStatus `gorm:"foreignKey:StatusId;references:Id"`
	Amount float64
}

// type PaymentStatus struct {
// 	Id int
// 	Name string
// }

// func (PaymentStatus) TableName() string {
//     return "payment_status"
// }