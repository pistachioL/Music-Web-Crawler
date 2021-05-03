package redis

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"time"
)

func Conn() *redis.Client{
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		//Password: "123456",
		//DB:       0,
	})
	//延迟到程序结束关闭链接
	//defer client.Close()
	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println("redis ping error", err.Error())
	}
	fmt.Println("redis ping success:", pong)
	return client
}

type RecentSong struct {
	title 	string	`json:"title"`
	author 	string	`json:"author"`
	url 	string	`json:"url"`
	pic		string	`json:"pic"`
	lrc 	string	`json:"lrc"`
}

var response = make([]string, 0)

func (m *RecentSong) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

/*
	@desc 以用户id作为key 每次点击播放时把歌曲列表放入redis。查看最近播放时，从Redis中取出一个月内播放过的歌曲。
	@return 返回最近播放的歌曲
 */
func SetRecentPlay(c *gin.Context) { //json->map->key
	client := Conn()
	key := c.Query("user") //用户名作为key
	recentSongList := c.Query("play") //歌曲json
	var recentSong RecentSong
	recentSong.MarshalBinary()
	setExpire(key)
	setList, err := client.SAdd(key, recentSongList).Result()
	if err != nil {
		log.Println("SAdd failed:", err)
		return
	}
	c.JSON(http.StatusOK, setList)
}


func GetRecentPlay(c *gin.Context) {
	client := Conn()
	key := c.Query("user") //用户名作为key
	songRes, _ :=  client.SMembers(key).Result() //["{}", "{}"]
	c.JSON(http.StatusOK, songRes)
}

func setExpire(key string) bool{
	client := Conn()
	res, err := client.Expire(key, 720*time.Hour).Result()
	if err != nil {
		fmt.Println("过期时间设置错误：",err)
	}
	fmt.Print("设置为期一个月的过期时间：",res)
	return res
}

func savePopoularList() {

}






