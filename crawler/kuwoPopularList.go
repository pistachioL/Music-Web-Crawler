package crawler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

type Kuwo struct {
	Img string`json:"img"`
	SongName 	string 	`json:"song_name"`
	
}

func KuwoSongList(context *gin.Context) {
	client := &http.Client{}
	url := "http://www.kuwo.cn/api/www/bang/bang/musicList?bangId=93&pn=1&rn=30&reqId=45fb0260-adc6-11eb-bbb1-97d257ecb717"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("请求酷我音乐错误",err)
		return
	}

	//设置请求头
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.109 Safari/537.36")
	request.Header.Set("Cookie", "GA1.2.397103415.1615459285; Hm_lvt_cdb524f42f0ce19b169a8071123a4797=1618463532,1619067810,1619840806,1620146654; _gid=GA1.2.15798403.1620146665; Hm_lpvt_cdb524f42f0ce19b169a8071123a4797=1620235260; _gat=1; kw_token=7I97Y4087OC")
	request.Header.Set("CSRF", "7I97Y4087OC")

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("请求出错！")
	}
	defer resp.Body.Close()
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(r))
	context.JSON(http.StatusOK, string(r)) //返回给前端
}
