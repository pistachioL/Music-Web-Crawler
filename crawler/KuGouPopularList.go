package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/jaydenwen123/go-util"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//func crawlerSongName(html string) {
//	doc, err := goquery.NewDocument(html)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	doc.Find("a[class=pc_temp_songname]").Each(func(i int, selection *goquery.Selection) {
//		selection.Attr("href")
//		res:= selection.Text()
//		fmt.Println(res)
//	})
//}
//
//func crawlerSongLink(html string) {
//	doc, err := goquery.NewDocument(html)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
//		doc.Find("a[class=pc_temp_songname]").Each(func(i int, selection *goquery.Selection) {
//			href, _ := selection.Attr("href")
//			fmt.Println(href)
//		})
//	})
//}


//通过get发送请求，返回数据
//第一个参数为字节数组，第二个参数为默认编码为utf-8的字符串
func RequestWithHeader(url string, headers map[string]string) ([]byte, string) {
	//1.发请求，获取数据
	//如果需要自己设置请求头，则通过http.NewRequest
	//resp, err := http.GetIndex(url)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logs.Error("RequestWithHeader http->NewRequest error:%v", err)
		return  nil,""
	}
	//设置请求头
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.109 Safari/537.36")
	if headers != nil {
		for key, value := range headers {
			request.Header.Add(key, value)
		}
	}
	//发送请求
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 120 * time.Second,
		}).Dial,
		IdleConnTimeout:       120 * time.Second,
		TLSHandshakeTimeout:   120 * time.Second,
		ResponseHeaderTimeout: 120 * time.Second,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   120 * time.Second,
	}
	resp, err := client.Do(request)
	if err != nil {
		logs.Error("http get error:", err.Error())
		//panic(err.Error())
		return nil, ""
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("ioutil ReadAll error:", err.Error())
		return nil, ""
	}
	if err = resp.Body.Close(); err != nil {
		logs.Error("resp Body Close error:", err.Error())
		return nil, ""
	}
	return content, string(content)
}

//解析出榜单的歌曲信息，返回的是json字符串
func ParseBoardSongsInfo(url string) string {
	_, data := RequestWithHeader(url, nil)
	//得到json数据
	data = util.MatchStringValue(`global.features =(?s:(.*?))}\]`, data)
	data = data + "}]"
	return data
}

func getHash(data string) ([]gjson.Result) {
	//获取hash值
	hashs := gjson.Get(data, "#.Hash")
	return hashs.Array()
}

func getAlbumId(data string) ([]gjson.Result) {
	//获取album_id值
	album_ids := gjson.Get(data, "#.album_id")
	return album_ids.Array()
}


func getSongRequestUrls(url string) []string {
	data := ParseBoardSongsInfo("https://www.kugou.com/yy/html/rank.html")
	hashs := getHash(data)
	album_ids := getAlbumId(data)

	//golang模拟请求
	var req, _ = http.NewRequest("GET", url, nil)
	q := req.URL.Query()
	songLen := len(hashs)
	sliec := make([]string, songLen)
	for v := range hashs {
		if(v != 0) {
			q.Del("hash")
			q.Del("album_id")
		}
		q.Add("hash", hashs[v].String())
		q.Add("album_id", album_ids[v].String())
		req.URL.RawQuery = q.Encode()
		sliec[v] = req.URL.String()
	}

	return sliec
}

//type NetEaseCloudMusicSong struct {
//	Status interface{} `json:"status"`
//	Errcode interface{} `json:"err_code"`
//	//Data Data `json:"data"`
//	Data *json.RawMessage `json:"data"`
//}
//
//type Data struct {
//	AudioName string `json:"audio_name"`
//	AlbumName string `json:"album_name"`
//	Image string `json:"img"`
//	SongName string `json:"song_name"`
//	AuthorName string `json:"author_name"`
//	Lyrics string `json:"lyrics"`
//	PlayUrl string `json:"play_url"`
//}

type Song struct {
	Status int `json:"status"`
	ErrCode int `json:"err_code"`
	Data Data `json:"data"`
}
type Authors struct {
	AuthorID string `json:"author_id"`
	AuthorName string `json:"author_name"`
	IsPublish string `json:"is_publish"`
	SizableAvatar string `json:"sizable_avatar"`
	Avatar string `json:"avatar"`
}
type Data struct {
	Hash string `json:"hash"`
	Timelength int `json:"timelength"`
	Filesize int `json:"filesize"`
	AudioName string `json:"audio_name"`
	HaveAlbum int `json:"have_album"`
	AlbumName string `json:"album_name"`
	AlbumID string `json:"album_id"`
	Img string `json:"img"`
	HaveMv int `json:"have_mv"`
	VideoID string `json:"video_id"`
	AuthorName string `json:"author_name"`
	SongName string `json:"song_name"`
	Lyrics string `json:"lyrics"`
	AuthorID string `json:"author_id"`
	Privilege int `json:"privilege"`
	Privilege2 string `json:"privilege2"`
	PlayURL string `json:"play_url"`
	Authors []Authors `json:"authors"`
	IsFreePart int `json:"is_free_part"`
	Bitrate int `json:"bitrate"`
	RecommendAlbumID string `json:"recommend_album_id"`
	AudioID string `json:"audio_id"`
	HasPrivilege bool `json:"has_privilege"`
	PlayBackupURL string `json:"play_backup_url"`
}

func unicode2utf8(source string) string {
	var res = []string{""}
	sUnicode := strings.Split(source, "\\u")
	var context = ""
	for _, v := range sUnicode {
		var additional = ""
		if len(v) < 1 {
			continue
		}
		if len(v) > 4 {
			rs := []rune(v)
			v = string(rs[:4])
			additional = string(rs[4:])
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			context += v
		}
		context += fmt.Sprintf("%c", temp)
		context += additional
	}
	res = append(res, context)
	return strings.Join(res, "")
}

func getSongDetails(url string) ([]Song,error) {
	urls := getSongRequestUrls(url)
	s := make([]Song, 0)
//	song := NetEaseCloudMusicSong{}
	var song Song
//	var data Data
//	var song json.RawMessage
	for i := range urls {
		resp, err := http.Get(urls[i])
		if err != nil {
			fmt.Print("请求飙升榜失败:", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Print("获取请求体失败:", err)
		}
		str := string(body)
		//	rex := regexp.MustCompile(`\(([^)]+)\)`) //匹配json
		rex := regexp.MustCompile(`\((.*)\)`) //匹配json
		out := rex.FindAllStringSubmatch(str, -1)

		//fmt.Println(out[0][1]) //去掉括号
		for _, i := range out {
			jsonData := []byte(i[1])
			//fmt.Print(unicode2utf8(string(jsonData)))
			err := json.Unmarshal(jsonData, &song) //反序列化为结构体
			if err != nil {
				fmt.Println("json解析失败：",err)
			}
			s = append(s, song) //把每个song存入map中
		}
	}

	return s, nil

}

func HandleSongData(context *gin.Context) {
	songDetails, err:= getSongDetails("https://wwwapi.kugou.com/yy/index.php?r=play/getdata&callback=jQuery191045751768061608544_1615257951217&dfid=3LjnlA1XAW9s3cB5ld2oVr1V&mid=99467f8a47af4fa16dc26fc68bab9215&platid=4&_=1615257951219")
	if err != nil {
		fmt.Println("获取歌曲详情信息错误getSongDetails()",err)
	}
	context.JSON(http.StatusOK, songDetails) //返回给前端
}

