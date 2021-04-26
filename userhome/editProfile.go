package userhome

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"net/http"
)

type User struct {
	Id		 int        `json:"id" `
	Username string	    `json:"username" db:"username"`
	Password string		`json:"password" db:password`
	Email    string		`json:"email" db:email`
	Gender   string		`json:"gender"`
	Avatar 	 string		`json:"avatar"`
	Desc 	 string		`json:"desc"`
}
var Response = make(map[string]interface{})
func conn() *gorm.DB{
	db,err := gorm.Open("mysql","root:971113Cg@@tcp(localhost)/music?charset=utf8&parseTime=True&loc=Local")
	if err != nil{
		fmt.Print("connect databases fail", err)
	}
	fmt.Print("connect database success")
	return db
}

func UpdateProfile(c *gin.Context) {
	db := conn()
	name := c.Request.FormValue("username")
	gender := c.Request.FormValue("gender")
	desc := c.Request.FormValue("desc")
	if len(name) == 0 && len(gender) == 0 && len(desc) == 0 {
		Response["code"] = -1
		Response["msg"] = "没有修改信息"
		return
	}
	var user = c.Query("user")

	//修改用户名
	if name != "" {
		db.Table("user").Where("username = ?", user).Update("username", name)
		Response["code"] = 0
		Response["msg"] = "修改成功"
	}

	//修改性别
	if gender != "" {
		db.Table("user").Where("username = ?", user).Update("gender", gender)
		Response["code"] = 0
		Response["msg"] = "修改成功"
	}

	//修改个人简介
	if desc  != "" {
		db.Table("user").Where("username = ?", user).Update("desc", desc)
		Response["code"] = 0
		Response["msg"] = "修改成功"
	}

	c.JSON(http.StatusOK, Response)
}
