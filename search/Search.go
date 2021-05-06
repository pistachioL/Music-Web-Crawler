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
	"githubLogin/model"
)


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
func handleSongDetail(searchApi string, c *gin.Context) model.SearchReq{
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
	r := model.SearchReq{}
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


func getSearchDetails(c *gin.Context) []model.DetailReq{
	urls := getSearchResUrls(c)
	searchSongs := make([]model.DetailReq, 0)
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
		searchRes := model.DetailReq{}
		body, _ := ioutil.ReadAll(req.Body) //可用string(body)查看，内容是多个json：{}{}
		json.Unmarshal(body, &searchRes)
		searchSongs = append(searchSongs, searchRes) //把每个song存入map中
	}
	//saveSearchRes(searchSongs)
	//如果es中有搜索结果，则返回es数据；否则重新爬取
	//getSearchResult(getKeyword(c))
	return searchSongs
}

func HandleSearch(context *gin.Context) {


	res := getSearchDetails(context)
	context.JSON(http.StatusOK, res) //返回给前端
}




