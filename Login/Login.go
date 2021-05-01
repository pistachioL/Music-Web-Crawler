package login

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"githubLogin/model"
	"net/http"
)


var UserInfo []model.User
var Response = make(map[string]interface{})


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
	db := model.Conn()
	var user model.User
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
	db := model.Conn()
	var user model.User
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
	db := model.Conn()
	newUser := &model.User{Username: name, Password: pwd, Email: email} // 根据指针找到数据表
	db.SingularTable(true)
	db.Create(newUser)
}




