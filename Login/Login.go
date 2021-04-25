package login

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
var UserInfo []User
var Response = make(map[string]interface{})

func conn() *gorm.DB{
	//cfg, err := ini.Load("conf/my.ini")
	//if err != nil {
	//	fmt.Printf("Fail to read configure file: %v", err)
	//	os.Exit(1)
	//}
	//gorm.Open("mysql", "conf/")
	//fmt.Print(cfg.Section("mysql").GetKey("User"))
	db,err := gorm.Open("mysql","root:971113Cg@@tcp(localhost)/music?charset=utf8&parseTime=True&loc=Local")
	if err != nil{
		fmt.Print("connect databases fail", err)
	}
	fmt.Print("connect database success")
	return db
}


//注册
func Register(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	email := c.Request.FormValue("email")
	//为空判断
	if len(username) == 0  || len(password) == 0|| len(email) == 0  {
		Response["code"] = -1
		Response["msg"] = "用户名或密码或邮箱不为空"
		c.JSON(http.StatusOK, Response)
		return
	}
	db := conn()
	var user User
	db.Table("user").Select("email").Where("email = ?", email).Scan(&user)

	if user.Email != "" {
		Response["code"] = 1
		Response["msg"] = "用户已存在"
	} else {
		AddUser(username, password, email)
		Response["code"] = 0
		Response["msg"] = "注册成功"
	}
	c.JSON(http.StatusOK, Response)
}

func Login(c *gin.Context) {
	name := c.Request.FormValue("username")
	pwd := c.Request.FormValue("password")
	if len(name) == 0 || len(pwd) == 0 {
		Response["code"] = -1
		Response["msg"] = "用户名密码不为空"
		c.JSON(http.StatusOK, Response)
		return
	}
	db := conn()
	var user User
	db.Table("user").Where("username = ?", name).Find(&user)
	if user.Email != "" {
		if user.Password == pwd {
			Response["code"] = 0
			Response["msg"] = "登录成功"
		} else{
			Response["code"] = 1
			Response["msg"] = "密码错误"
		}
	} else {
		Response["code"] = -1
		Response["msg"] = "登录失败，用户尚未注册"
	}
	c.JSON(http.StatusOK, Response)
}


func AddUser(name string, pwd string, email string) {
	db := conn()
	newUser := &User{Username: name, Password: pwd, Email: email} // 根据指针找到数据表
	db.SingularTable(true)
	db.Create(newUser)
}




