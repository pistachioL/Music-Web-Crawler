package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
	"github.com/jaydenwen123/go-util"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"time"
)

func crawlerSongName(html string) {
	doc, err := goquery.NewDocument(html)
	if err != nil {
		fmt.Println(err)
	}

	doc.Find("a[class=pc_temp_songname]").Each(func(i int, selection *goquery.Selection) {
		selection.Attr("href")
		res:= selection.Text()
		fmt.Println(res)
	})
}

func crawlerSongLink(html string) {
	doc, err := goquery.NewDocument(html)
	if err != nil {
		fmt.Println(err)
	}

	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		doc.Find("a[class=pc_temp_songname]").Each(func(i int, selection *goquery.Selection) {
			href, _ := selection.Attr("href")
			fmt.Println(href)
		})
	})
}


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
	sliec := make([]string, len(hashs))
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


func getSongDetails(url string) string {
	urls := getSongRequestUrls(url)
	//for i := range urls {

	//}
	resp, err := http.Get(urls[1])
	if err != nil {
		fmt.Print("请求飙升榜失败:", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	str := string(body)

	rex := regexp.MustCompile(`\((.*)\)`) //匹配json
	out := rex.FindAllStringSubmatch(str, -1)
	var jsonStr string = ""
	for _, i := range out {
		jsonStr += i[1]
	}
	fmt.Print(jsonStr)

	return jsonStr
}


type Song struct {
	Status interface{} `json:"status"`
	Errcode interface{} `json:"err_code"`
//	Data interface{} `json:"data"`
	Data Data
}
type Data struct {
	AudioName string `json:"audio_name"`
	AlbumName string `json:"album_name"`
}


func main() {
	//getSongDetails("https://wwwapi.kugou.com/yy/index.php?r=play/getdata&callback=jQuery191045751768061608544_1615257951217&dfid=3LjnlA1XAW9s3cB5ld2oVr1V&mid=99467f8a47af4fa16dc26fc68bab9215&platid=4&_=1615257951219")

	data := `{"status":1,"err_code":0,"data":{"hash":"AD70C2B0D2B213E31960D1FB4AE6A002","timelength":248087,"filesize":3982359,"audio_name":"BEYOND - \u6211\u662f\u6124\u6012","have_album":1,"album_name":"\u4e50\u4e0e\u6012","album_id":"973001","img":"http:\/\/imge.kugou.com\/stdmusic\/20150715\/20150715232800432202.jpg","have_mv":1,"video_id":"598588","author_name":"BEYOND","song_name":"\u6211\u662f\u6124\u6012","lyrics":"\ufeff[id:$00000000]\r\n[ar:beyond]\r\n[ti:\u6211\u662f\u6124\u6012]\r\n[by:]\r\n[hash:15749c637fd980a5c3c201d3994f489c]\r\n[al:]\r\n[sign:]\r\n[qq:]\r\n[total:252160]\r\n[offset:0]\r\n[00:01.80]\u4f5c\u8bcd\uff1a\u9ec4\u8d2f\u4e2d\r\n[00:02.96]\u4f5c\u66f2\uff1a\u9ec4\u5bb6\u9a79\r\n[00:15.29]Woo AI\r\n[00:18.15]\u53ef\u5426\u4e89\u756a\u4e00\u56d7\u6c14\r\n[00:32.11]\u6211\u662f\u6076\u68a6\r\n[00:35.08]\u5929\u5929\u90fd\u53ef\u9a9a\u6270\u4f60\r\n[00:38.24]\u4e0e\u4f60\u9047\u7740\u5728\u8def\u9014\r\n[00:41.30]\u4f60\u83ab\u9000\u907f\r\n[00:44.50]\u6211\u662f\u6124\u6012\r\n[00:47.56]\u5206\u5206\u949f\u53ef\u70e7\u6b7b\u4f60\r\n[00:50.72]\u51e0\u591a\u865a\u5047\u7684\u597d\u6c49\r\n[00:53.88]\u90fd\u7747\u4e0d\u8d77\r\n[00:56.76]\u53ea\u60f3\u541e\u5343\u5428\u7684\u6012\u706b\r\n[01:00.16]\u672a\u53bb\u60f3\u5931\u58f0\u547c\u53eb\r\n[01:03.22]I'll never die\r\n[01:04.82]I'll never cry\r\n[01:06.58]you'll see\r\n[01:10.09]WOO AI\r\n[01:13.10]\u53ef\u5426\u4e89\u756a\u4e00\u56d7\u6c14\r\n[01:16.21]WOO AI\r\n[01:19.43]\u771f\u672c\u6027\r\n[01:21.09]\u600e\u53ef\u4ee5\u6539\r\n[01:31.74]\u4f60\u52ff\u8bf4\u8bdd\r\n[01:34.80]\u7686\u56e0\u4eca\u5929\u7684\u771f\u7406\r\n[01:37.84]\u8bb2\u8d77\u59cb\u7ec8\u90fd\u8ddf\u6211\r\n[01:40.90]\u6709\u6bb5\u8ddd\u79bb\r\n[01:44.11]\u62d2\u7edd\u5bf9\u8bdd\r\n[01:47.22]\u7686\u56e0\u4eca\u5929\u7684\u5929\u6c14\r\n[01:50.38]\u600e\u6837\u547c\u5438\u90fd\u4e0d\u60ef\r\n[01:53.48]\u592a\u6ca1\u8da3\u5473\r\n[01:56.37]\u53ea\u60f3\u541e\u5343\u5428\u7684\u6012\u706b\r\n[01:59.69]\u672a\u53bb\u60f3\u5931\u58f0\u547c\u53eb\r\n[02:02.77]I'll never die\r\n[02:04.21]I'll never cry\r\n[02:06.12]you'll see\r\n[02:09.63]WOO AI\r\n[02:12.69]\u53ef\u5426\u4e89\u756a\u4e00\u56d7\u6c14\r\n[02:16.15]WOO AI\r\n[02:18.96]\u771f\u672c\u6027\r\n[02:20.74]\u600e\u53ef\u4ee5\u6539\r\n[02:24.83]Come on\r\n[02:56.12]\u53ea\u60f3\u541e\u5343\u5428\u7684\u6012\u706b\r\n[02:59.38]\u672a\u53bb\u60f3\u5931\u58f0\u547c\u53eb\r\n[03:02.38]I'll never die\r\n[03:03.93]I'll never cry\r\n[03:05.73]you'll see\r\n[03:09.34]WOO AI\r\n[03:12.30]\u53ef\u5426\u4e89\u756a\u4e00\u56d7\u6c14\r\n[03:15.42]WOO AI\r\n[03:18.58]\u771f\u672c\u6027 \u600e\u53ef\u4ee5\u6539\r\n[03:24.85]WOO AI\r\n[03:27.95]\u53ef\u5426\u4e89\u756a\u4e00\u56d7\u6c14\r\n[03:31.11]WOO AI\r\n[03:34.22]\u771f\u672c\u6027 \u600e\u53ef\u4ee5\u6539\r\n","author_id":"7249","privilege":10,"privilege2":"1010","play_url":"https:\/\/webfs.yun.kugou.com\/202103111218\/a5a965468f2acec3c9d15d784403881c\/part\/0\/960931\/G141\/M09\/04\/17\/LYcBAFuGce-ARlWtADzEF7I63RQ112.mp3","authors":[{"author_id":"7249","author_name":"BEYOND","is_publish":"1","sizable_avatar":"http:\/\/singerimg.kugou.com\/uploadpic\/softhead\/{size}\/20160418\/20160418100531196.jpg","avatar":"http:\/\/singerimg.kugou.com\/uploadpic\/softhead\/400\/20160418\/20160418100531196.jpg"}],"is_free_part":1,"bitrate":128,"recommend_album_id":"973001","audio_id":"261918","has_privilege":true,"play_backup_url":"https:\/\/webfs.cloud.kugou.com\/202103111218\/24a45c533dbaafeeaa03e07dd39b1de7\/part\/0\/960931\/G141\/M09\/04\/17\/LYcBAFuGce-ARlWtADzEF7I63RQ112.mp3","trans_param":{"hash_offset":{"start_byte":0,"end_byte":960931,"start_ms":0,"end_ms":60000,"offset_hash":"99F3DEBDE241F06ED38EEF2C92B92CED","file_type":0},"musicpack_advance":1,"pay_block_tpl":1,"display":32,"display_rate":1,"cpy_grade":5,"cpy_level":1,"cid":4748837,"cpy_attr0":0}}}`
	//data := `{"status":1,"err_code":0, "data":{"audio_name":"BEYOND - \u6211\u662f\u6124\u6012"}}`
	str := []byte(data)
	song := Song{}
	err := json.Unmarshal(str, &song)
	if err != nil {
		fmt.Println("json解析失败：",err)
	}
	fmt.Print(song)

	//	crawlerSongName("https://www.kugou.com/yy/html/rank.html")
	//	ParseBoardSongsInfo("https://www.kugou.com/yy/html/rank.html")

}
