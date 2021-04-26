package userhome

import (
	"github.com/gin-gonic/gin"
	"net/http"
)
var readResponse = make(map[string]interface{})
func ReadProfile(c *gin.Context) {

	var user User
	db := conn()
	var currentUser = c.Query("user")
	db.Table("user").Where("username = ?", currentUser).Find(&user)
	if user.Username != "" {
		readResponse["code"] = 0
		readResponse["msg"] = user
	} else {
		readResponse["code"] = -1
		readResponse["msg"] = "没有该用户"
	}
	c.JSON(http.StatusOK, readResponse)
}
