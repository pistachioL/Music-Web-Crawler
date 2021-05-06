package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)
type Req struct {
	Code int `json:"code"`
	CurTime int64 `json:"curTime"`
	Data Data `json:"data"`
	Msg string `json:"msg"`
	ProfileID string `json:"profileId"`
	ReqID string `json:"reqId"`
	TID string `json:"tId"`
}
type Mvpayinfo struct {
	Play int `json:"play"`
	Vid int `json:"vid"`
	Down int `json:"down"`
}
type FeeType struct {
	Song string `json:"song"`
	Vip string `json:"vip"`
}
type PayInfo struct {
	Play string `json:"play"`
	Download string `json:"download"`
	LocalEncrypt string `json:"local_encrypt"`
	Limitfree int `json:"limitfree"`
	CannotDownload int `json:"cannotDownload"`
	ListenFragment string `json:"listen_fragment"`
	CannotOnlinePlay int `json:"cannotOnlinePlay"`
	FeeType FeeType `json:"feeType"`
	Down string `json:"down"`
}
type Data struct {
	Img string `json:"img"`
	Num string `json:"num"`
	Pub string `json:"pub"`
	MusicList []MusicList `json:"musicList"`
}

type MusicList struct {
	Musicrid string `json:"musicrid"`
	Barrage string `json:"barrage"`
	Artist string `json:"artist"`
	Trend string `json:"trend"`
	Pic string `json:"pic"`
	Isstar int `json:"isstar"`
	Rid int `json:"rid"`
	Duration int `json:"duration"`
	Score100 string `json:"score100"`
	ContentType string `json:"content_type"`
	RankChange string `json:"rank_change"`
	Track int `json:"track"`
	HasLossless bool `json:"hasLossless"`
	Hasmv int `json:"hasmv"`
	ReleaseDate string `json:"releaseDate"`
	Album string `json:"album"`
	Albumid int `json:"albumid"`
	Pay string `json:"pay"`
	Artistid int `json:"artistid"`
	Albumpic string `json:"albumpic"`
	Originalsongtype int `json:"originalsongtype"`
	SongTimeMinutes string `json:"songTimeMinutes"`
	IsListenFee bool `json:"isListenFee"`
	Pic120 string `json:"pic120"`
	Name string `json:"name"`
	Online int `json:"online"`
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
	res := string(r)
	var kuwo Req
	if err := json.Unmarshal([]byte(res), &kuwo); err != nil {
		fmt.Println("json反序列化错误", err)
	}
	context.JSON(http.StatusOK, kuwo) //返回给前端
}
