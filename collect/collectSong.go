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
	defer db.Close()
	db.LogMode(true)
	db.AutoMigrate(&model.User{}, &model.Song{}) //生成中间表
	var currentUser = c.Query("user")
	var user model.User
	var songId = c.Query("songId")
	var song []model.Song

	db.Model(model.Song{}).Where("id = ?" , songId).Find(&song)
	userId := db.Select("id").Where("username = ?", currentUser).Find(&user)


	//songId, err := strconv.Atoi(id)
	//if err != nil {
	//	fmt.Println("songId类型转换错误", err)
	//}
	db.Table("users").Where("id = ?", userId).First(&user) //查找点击收藏的歌曲
	//db.Model(&user).Association("Song").Find(&song)
	res := db.Model(&user).Association("Song").Append(&song) //建立歌曲和用户的关联

	fmt.Println("用户歌曲关联：", res)
	//user1 := model.User{Username: currentUser, Song: song}
	//db.Save(user1)


	//db.Model(&user).Association("Song")

	//addCollect := &model.UserSong{Username: currentUser, Songid: songId}
	//var user model.User
	//addCollect :=
	//res := db.Create(&addCollect)
	//fmt.Print("收藏结果：",res)
	c.JSON(http.StatusOK, res) //返回给前端
}


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

//func CancelCollection(c *gin.Context) {
//	db := model.Conn()
//	db.SingularTable(true)
//	var currentUser = c.Query("user")
//}

