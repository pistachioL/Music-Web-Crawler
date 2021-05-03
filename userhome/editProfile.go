package userhome

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"githubLogin/model"
	"net/http"
)



var Response = make(map[string]interface{})


func UpdateProfile(c *gin.Context) {
	db := model.Conn()
	defer db.Close()
	var user model.User
	var currentUser = c.Query("user")
	name := c.Request.FormValue("username")
	gender := c.Request.FormValue("gender")
	desc := c.Request.FormValue("desc")

	db.Table("users").Where("username = ?", currentUser).Find(&user)
	if user.Username != "" { //有这个用户
		Response["code"] = ""
		Response["msg"] = ""
		Response["gender"] = ""
		Response["name"] = ""
		Response["desc"] = ""
		if user.Username == name && user.Gender == gender && user.Desc == desc {
			db.Table("user").Where("username = ?", currentUser).Update("username", name)
			Response["code"] = -2
			Response["msg"] = "没有修改的信息"
			c.JSON(http.StatusOK, Response)
			return
		} else {
			//修改用户名
			if name != user.Username {
				db.Table("users").Where("username = ?", currentUser).Update("username", name)
				Response["code"] = 0
				Response["name"] = "用户名修改成功"
			}

			//修改性别
			if gender != user.Gender {
				db.Table("users").Where("username = ?", currentUser).Update("gender", gender)
				Response["code"] = 0
				Response["gender"] = "性别修改成功"
			}

			//修改个人简介
			if desc != user.Desc {
				db.Table("users").Where("username = ?", currentUser).Update("desc", desc)
				Response["code"] = 0
				Response["desc"] = "个人简介修改成功"
			}

			c.JSON(http.StatusOK, Response)
			return
		}
	} else {
		Response["code"] = -1
		Response["msg"] = "没有该用户"
		c.JSON(http.StatusOK, Response)
		return
	}
}
