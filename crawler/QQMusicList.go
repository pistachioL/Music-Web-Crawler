package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)
type QQReq struct {
	Code int `json:"code"`
	Color int `json:"color"`
	CommentNum int `json:"comment_num"`
	CurSongNum int `json:"cur_song_num"`
	Date string `json:"date"`
	DayOfYear string `json:"day_of_year"`
	SongBegin int `json:"song_begin"`
	Songlist []Songlist `json:"songlist"`
	Topinfo Topinfo `json:"topinfo"`
	TotalSongNum int `json:"total_song_num"`
	UpdateTime string `json:"update_time"`
}
type Pay struct {
	Payalbum int `json:"payalbum"`
	Payalbumprice int `json:"payalbumprice"`
	Paydownload int `json:"paydownload"`
	Payinfo int `json:"payinfo"`
	Payplay int `json:"payplay"`
	Paytrackmouth int `json:"paytrackmouth"`
	Paytrackprice int `json:"paytrackprice"`
	Timefree int `json:"timefree"`
}
type Preview struct {
	Trybegin int `json:"trybegin"`
	Tryend int `json:"tryend"`
	Trysize int `json:"trysize"`
}
type Singer struct {
	ID int `json:"id"`
	Mid string `json:"mid"`
	Name string `json:"name"`
}
type QQData struct {
	Albumdesc string `json:"albumdesc"`
	Albumid int `json:"albumid"`
	Albummid string `json:"albummid"`
	Albumname string `json:"albumname"`
	Alertid int `json:"alertid"`
	BelongCD int `json:"belongCD"`
	CdIdx int `json:"cdIdx"`
	Icons int `json:"icons"`
	Interval int `json:"interval"`
	Isonly int `json:"isonly"`
	Label string `json:"label"`
	Msgid int `json:"msgid"`
	Pay Pay `json:"pay"`
	Preview Preview `json:"preview"`
	Rate int `json:"rate"`
	Singer []Singer `json:"singer"`
	Size128 int `json:"size128"`
	Size320 int `json:"size320"`
	Size51 int `json:"size5_1"`
	Sizeape int `json:"sizeape"`
	Sizeflac int `json:"sizeflac"`
	Sizeogg int `json:"sizeogg"`
	Songid int `json:"songid"`
	Songmid string `json:"songmid"`
	Songname string `json:"songname"`
	Songorig string `json:"songorig"`
	Songtype int `json:"songtype"`
	StrMediaMid string `json:"strMediaMid"`
	Stream int `json:"stream"`
	Switch int `json:"switch"`
	Type int `json:"type"`
	Vid string `json:"vid"`
}
type Songlist struct {
	FrankingValue string `json:"Franking_value"`
	CurCount string `json:"cur_count"`
	QQData QQData `json:"data"`
	InCount string `json:"in_count"`
	OldCount string `json:"old_count"`
}
type Topinfo struct {
	ListName string `json:"ListName"`
	MacDetailPicURL string `json:"MacDetailPicUrl"`
	MacListPicURL string `json:"MacListPicUrl"`
	UpdateType string `json:"UpdateType"`
	Albuminfo string `json:"albuminfo"`
	HeadPicV12 string `json:"headPic_v12"`
	Info string `json:"info"`
	Listennum int `json:"listennum"`
	Pic string `json:"pic"`
	PicDetail string `json:"picDetail"`
	PicAlbum string `json:"pic_album"`
	PicH5 string `json:"pic_h5"`
	PicV11 string `json:"pic_v11"`
	PicV12 string `json:"pic_v12"`
	TopID string `json:"topID"`
	Type string `json:"type"`
}

func QQSongList(context *gin.Context) {
	url := "https://c.y.qq.com/v8/fcg-bin/fcg_v8_toplist_cp.fcg?g_tk=5381&uin=0&format=json&inCharset=utf-8&outCharset=utf-8¬ice=0&platform=h5&needNewCode=1&tpl=3&page=detail&type=top&topid=27&_=1519963122923"
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("请求酷我音乐错误",err)
		return
	}
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
	var qqMusic QQReq
	if err := json.Unmarshal([]byte(res), &qqMusic); err != nil {
		fmt.Println("qq音乐反序列化错误", err)
		return
	}

	context.JSON(http.StatusOK, qqMusic) //返回给前端
}


