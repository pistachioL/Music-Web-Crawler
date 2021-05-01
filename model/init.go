package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)


var db *gorm.DB

func Conn() *gorm.DB{
	db,err := gorm.Open("mysql","root:971113Cg@@tcp(localhost)/music?charset=utf8&parseTime=True&loc=Local")
	if err != nil{
		fmt.Print("connect databases fail", err)
	}
	fmt.Print("connect database success")
	return db
}
