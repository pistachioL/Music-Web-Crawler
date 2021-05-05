package collect

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"githubLogin/model"
	"net/http"
)


func AddCollection(c *gin.Context) {
	db := model.Conn()
	//db.LogMode(true)
	db.AutoMigrate(&model.User{}, &model.Song{}) //生成中间表
	var currentUser = c.Query("user")
	var user model.User
	var songId = c.Query("songId")
	var song []model.Song
	db.Model(model.Song{}).Where("id = ?" , songId).Find(&song)
	userId := db.Select("id").Where("username = ?", currentUser).Find(&user)
	db.Table("users").Where("id = ?", userId).First(&user) //查找点击收藏的歌曲
	res := db.Model(&user).Association("Song").Append(&song) //建立歌曲和用户的关联

	fmt.Println("用户歌曲关联：", res)
	c.JSON(http.StatusOK, res) //返回给前端
}

//查询所有收藏记录
func QueryCollection(c *gin.Context) {
	db := model.Conn()
	defer db.Close()
	db.LogMode(true)
	var currentUser = c.Query("user")
	var user model.User
	var songs []model.Song
	db.Model(model.User{}).Where("username = ?", currentUser).Find(&user)
	db.Model(&user).Association("Song").Find(&songs)
	c.JSON(http.StatusOK, songs) //返回给前端
}

//只查询5条收藏记录
func QueryFiveCollection(c *gin.Context) {
	db := model.Conn()
	defer db.Close()
	db.LogMode(true)
	var currentUser = c.Query("user")
	var user model.User
	var songs []model.Song
	db.Model(model.User{}).Where("username = ?", currentUser).Find(&user)
	db.Model(&user).Association("Song").Find(&songs)
	c.JSON(http.StatusOK, songs) //返回给前端
}

//取消收藏
func CancelCollection(c *gin.Context) {
	db := model.Conn()
	defer db.Close()
	db.LogMode(true)
	//db.AutoMigrate(&model.User{}, &model.Song{}) //生成中间表
	var currentUser = c.Query("user")
	var user model.User
	var songId = c.Query("songId")
	var song []model.Song
	db.Model(model.Song{}).Where("id = ?" , songId).Find(&song)
	userId := db.Select("id").Where("username = ?", currentUser).Find(&user)
	db.Table("users").Where("id = ?", userId).First(&user) //查找取消收藏的歌曲
	res := db.Model(&user).Association("Song").Delete(&song) //建立歌曲和用户的关联

	fmt.Println("取消收藏：", res)
	c.JSON(http.StatusOK, res) //返回给前端
}

