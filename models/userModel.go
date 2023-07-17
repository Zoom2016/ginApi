package models

type User struct {
	Id       int    `json:"id" gorm:"column:id;primaryKey"`
	Name     string `json:"name" gorm:"column:name;default:null"`
	UserName string `json:"username" gorm:"column:username;default:null"`
	Password string `json:"password" gorm:"column:password;default:null"`
	Age      int    `json:"age" gorm:"column:age;default:null"`
	Cdt      XTime  `json:"cdt"`
	Mdt      XTime  `json:"mdt"`
}

func (User) TableName() string {
	return "user"
}
