package search

/*
 @desc 酷狗歌手或歌曲搜索
 @date 2021/4/15
*/

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)
type SearchReq struct {
	Status int `json:"status"`
	Error string `json:"error"`
	Data SearchData `json:"data"`
	Errcode int `json:"errcode"`
}

type SearchData struct {
	SearchInfo []SearchInfo `json:"info"`
}

type SearchInfo struct {
	Singername string `json:"singername"`
	Songname string `json:"songname"`
	Hash string `json:"hash"`
	AlbumAudioID int `json:"album_audio_id"`
	AlbumID string `json:"album_id"`
}

type DetailReq struct {
	Status int `json:"status"`
	ErrCode int `json:"err_code"`
	Data Data `json:"data"`
}
//type Authors struct {
//	AuthorID string `json:"author_id"`
//	AuthorName string `json:"author_name"`
//	IsPublish string `json:"is_publish"`
//	SizableAvatar string `json:"sizable_avatar"`
//	Avatar string `json:"avatar"`
//}
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
	//Authors []Authors `json:"authors"`
	IsFreePart int `json:"is_free_part"`
	Bitrate int `json:"bitrate"`
	RecommendAlbumID string `json:"recommend_album_id"`
	AudioID string `json:"audio_id"`
	HasPrivilege bool `json:"has_privilege"`
	PlayBackupURL string `json:"play_backup_url"`
}

//酷狗搜索api
var searchApi = "http://msearchcdn.kugou.com/api/v3/search/song?tagtype=全部&pagesize=50"
//酷狗详情api
var songDetailApi  = "https://wwwapi.kugou.com/yy/index.php?r=play/getdata&dfid=3LjnlA1XAW9s3cB5ld2oVr1V&mid=99467f8a47af4fa16dc26fc68bab9215&platid=4&_=1615257951219"

//获取输入框中的关键字
func getKeyword(c *gin.Context) string {
	var keyword = c.Query("keyword")
	fmt.Print(keyword)
	return keyword
}

/*
 @desc 请求搜索api中所有歌曲
 @param searchApi 搜索api
*/
func handleSongDetail(searchApi string, c *gin.Context) SearchReq {
	client, err := http.NewRequest(http.MethodGet, searchApi, nil)
	keyword := getKeyword(c)
	if err != nil {
		panic("请求失败")
	}
	params := client.URL.Query()
	params.Add("keyword", keyword)
	client.URL.RawQuery = params.Encode()
	req, err := http.DefaultClient.Do(client)
	defer func() {
		req.Body.Close()
	}()
	if req.StatusCode != 200 {
		panic("返回错误响应")
	}
	body, _ := ioutil.ReadAll(req.Body)
	r := SearchReq{}
	json.Unmarshal(body, &r)
	if r.Status != 1 {
		panic(r.Error)
	}
	return r
}

/*
 @des 获取搜索结果的api（各个api携带唯一的hash和album_id）
 @return 返回搜索结果
*/
func getSearchResUrls(c *gin.Context) []string{
	res := handleSongDetail(searchApi, c)
	songs := res.Data.SearchInfo
	searchLen := len(songs)
	hashs := make([]string, searchLen)
	albumids :=  make([]string, searchLen)
	searchRes := make([]string, searchLen)
	for k,song := range songs {
		hashs[k] = song.Hash
		albumids[k] = song.AlbumID
	}

	var req, _ = http.NewRequest("GET", songDetailApi, nil)
	q := req.URL.Query()
	for k := range hashs {
		if(k != 0) {
			q.Del("hash")
			q.Del("album_id")
		}
		q.Add("hash",hashs[k])
		q.Add("album_id", albumids[k])
		req.URL.RawQuery = q.Encode()
		searchRes[k] = req.URL.String()
	}
	return searchRes
}


func getSearchDetails(c *gin.Context) []DetailReq{
	urls := getSearchResUrls(c)
	searchSongs := make([]DetailReq, 0)
	for i := range urls {
		req, err := http.Get(urls[i])
		if err != nil {
			fmt.Print("请求搜索接口失败:", err)
		}
		defer func() {
			req.Body.Close()
		}()
		if req.StatusCode != 200 {
			panic("返回错误响应")
		}
		searchRes := DetailReq{}
		body, _ := ioutil.ReadAll(req.Body) //可用string(body)查看，内容是多个json：{}{}
		json.Unmarshal(body, &searchRes)
		searchSongs = append(searchSongs, searchRes) //把每个song存入map中
	}
	saveSearchRes(searchSongs)
	getSearchResult()
	return searchSongs
}

func HandleSearch(context *gin.Context) {
	res := getSearchDetails(context)
	context.JSON(http.StatusOK, res) //返回给前端
}




