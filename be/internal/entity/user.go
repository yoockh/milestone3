package entity

type Users struct {
	Id int
	Name string
	Email string
	Password string `json:"-"`
	Role string
	// Role Role `gorm:"foreignKey:RoleId;references:Id"`
}

// type Role struct {
// 	Id int
// 	Name string
// }