package collect

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)
type Collection struct {
	gorm.Model
	Id		 				int        	`json:"id" `
//	CollectUid 				string	    `json:"collect_uid" db:"collect_uid"`
	CollectSongTimelength 	string		`json:"collect_song_timelength" db:collect_song_timelength`
	CollectSongArt    		string		`json:"collect_song_art" db:collect_song_art`
	CollectAlbum  	 		string		`json:"collect_album"`
	CollectAlbumImg 	 	string		`json:"collect_album_img"`
	CollectSongSrc 			string		`json:"collect_song_src"`
	Users 					[]User 		`gorm:"many2many:user_collections;"`
}

type User struct {
	gorm.Model
	Id		 int        `json:"id" `
	Username string	    `json:"username" db:"username"`
	Password string		`json:"password" db:password`
	Email    string		`json:"email" db:email`
	Gender   string		`json:"gender"`
	Avatar 	 string		`json:"avatar"`
	Desc 	 string		`json:"desc"`
	Collections []Collection `gorm:"many2many:user_collections;"`
}

func conn()*gorm.DB {
	db,err := gorm.Open("mysql","root:971113Cg@@tcp(localhost)/music?charset=utf8&parseTime=True&loc=Local")
	if err != nil{
		fmt.Print("connect databases fail", err)
	}
	fmt.Print("connect database success")
	return db
}

func CollectSong(c *gin.Context) {
	db := conn()
	defer db.Close()
	var currentUser = c.Query("user")
	var user User
	db.Table("user").Where("username = ?", currentUser).Find(&user)

}


