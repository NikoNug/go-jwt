package models

type User struct {
	Id       int    `json:"id" gorm:"primaryKey"`
	FullName string `json:"full_name" gorm:"varchar(300)"`
	Username string `json:"username" gorm:"varchar(300)"`
	Password string `json:"password" gorm:"varchar(300)"`
}
