package main

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	crawlerSongInfo()
	//downloadMusic()

}

func getSongInfo(index int) error {
	musicInfoUrl := "http://player.kuwo.cn/webmusic/st/getNewMuiseByRid?rid=MUSIC_" + strconv.Itoa(index)
	log.Println("musicInfoUrl:", musicInfoUrl)
	res, err := http.Get(musicInfoUrl)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	var song Song
	var responseBytes []byte
	responseBytes, _ = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if len(responseBytes) > 60 { //song不为空
		responseStr := string(responseBytes)
		tempResponseStr := strings.Replace(responseStr, "&", "&amp;", -1) //xml解析前，将"&"替换为"&amp;"，否则解析会失败
		err = xml.Unmarshal([]byte(tempResponseStr), &song)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		if song.Music_id != "" {
			requestMp3Url := "http://antiserver.kuwo.cn/anti.s?rid=" + song.Music_id + "&format=mp3&type=convert_url&response=url"
			requestAacUrl := "http://antiserver.kuwo.cn/anti.s?rid=" + song.Music_id + "&format=aac&type=convert_url&response=url"
			log.Println("requestMp3Url", requestMp3Url)
			log.Println("requestAacUrl", requestAacUrl)
			mp3UrlRes, err := http.Get(requestMp3Url)
			if err != nil {
				log.Println(err.Error())
				return err
			}
			var mp3UrlResBytes []byte
			mp3UrlResBytes, _ = ioutil.ReadAll(mp3UrlRes.Body)
			defer mp3UrlRes.Body.Close()
			mp3Url := string(mp3UrlResBytes)
			log.Println("mp3Url:", mp3Url)
			aacUrlRes, err := http.Get(requestAacUrl)
			if err != nil {
				log.Println(err.Error())
				return err
			}
			var aacUrlResBytes []byte
			aacUrlResBytes, _ = ioutil.ReadAll(aacUrlRes.Body)
			defer aacUrlRes.Body.Close()
			aacUrl := string(aacUrlResBytes)
			log.Println("aacUrl:", aacUrl)
			if song.Song_url != "" { //去掉前缀
				song.Song_url = song.Song_url[21:len(song.Song_url)]
			}
			if song.Artist_url != "" { //去掉前缀
				song.Artist_url = song.Artist_url[21:len(song.Artist_url)]
			}
			//
			//err = insertIntoDb(&song, requestMp3Url, requestAacUrl, mp3Url, aacUrl)
			//if err != nil {
			//	log.Println(err.Error())
			//}
		}
	}

	return nil
}

type Song struct {
	XMLName       xml.Name `xml:"Song"`
	Music_id      string   `xml:"music_id"`
	Mv_rid        string   `xml:"mv_rid"`
	Name          string   `xml:"name"`
	Song_url      string   `xml:"song_url"`
	Artist        string   `xml:"artist"`
	Artid         string   `xml:"artid"`
	Singer        string   `xml:"singer"`
	Special       string   `xml:"special"`
	Ridmd591      string   `xml:"ridmd591"`
	Mp3size       string   `xml:"mp3size"`
	Artist_url    string   `xml:"artist_url"`
	Auther_url    string   `xml:"auther_url"`
	Playid        string   `xml:"playid"`
	Artist_pic    string   `xml:"artist_pic"`
	Artist_pic240 string   `xml:"artist_pic240"`
	Path          string   `xml:"path"`
	Mp3path       string   `xml:"mp3path"`
	Aacpath       string   `xml:"aacpath"`
	Wmadl         string   `xml:"wmadl"`
	Mp3dl         string   `xml:"mp3dl"`
	Aacdl         string   `xml:"aacdl"`
	Lyric         string   `xml:"lyric"`
	Lyric_zz      string   `xml:"lyric_zz"`
	RequestMp3Url string
	ReqeustAacUrl string
}

var dbConfig = "root:123456@tcp(localhost:3306)/music"
var db *sql.DB
var openDbError error

func init() {
	db, openDbError = sql.Open("mysql", dbConfig)
	if openDbError != nil {
		log.Println("open mysql failed,error:" + openDbError.Error())
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)
}
func insertIntoDb(song *Song, requestMp3Url, requestAacUrl, mp3Url, aacUrl string) error {
	if db != nil && song != nil {
		stmt, err := db.Prepare("insert into songinfo(music_id,mv_rid,name,song_url,artist,artid,singer,special,ridmd591,mp3size,artist_url,auther_url,playid,artist_pic,artist_pic240,path,mp3path,aacpath,wmadl,mp3dl,aacdl,lyric,lyric_zz,request_mp3_url,request_aac_url,song_mp3_url,song_aac_url)values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
		if err != nil {
			return err
		}
		rs, err := stmt.Exec(song.Music_id, song.Mv_rid, song.Name, song.Song_url, song.Artist, song.Artid, song.Singer, song.Special, song.Ridmd591, song.Mp3size, song.Artist_url, song.Auther_url, song.Playid, song.Artist_pic, song.Artist_pic240, song.Path, song.Mp3path, song.Aacpath, song.Wmadl, song.Mp3dl, song.Aacdl, song.Lyric, song.Lyric_zz, requestMp3Url, requestAacUrl, mp3Url, aacUrl)
		if err != nil {
			return err
		}
		//我们可以获得插入的id
		id, err := rs.LastInsertId()
		if err != nil {
			return err
		}
		log.Println("LastInsertId():", id)
		//可以获得影响行数
		affect, err := rs.RowsAffected()
		if err != nil {
			return err
		}
		log.Println("RowsAffected():", affect)
		return nil

	} else {
		return errors.New("db or song is nil")
	}
	return nil

}

//爬取歌曲信息，插入数据库
func crawlerSongInfo() {
	/**
	1、遍历所有的ID,拼接http://player.kuwo.cn/webmusic/st/getNewMuiseByRid?rid=MUSIC_ID
	2、向每个URL发送请求，解析，如果<Song></Song>不为空的话，将歌曲信息入库
	3、通过1和2应该可以拿到所有的音乐信息；拿个每一项歌曲信息后，拼接http://antiserver.kuwo.cn/anti.s?rid=MUSIC_ID&format=aac|mp3&type=convert_url&response=url
	4、向3中的连接发请求，拿到AAC和mp3的歌曲地址
	5、从100000开始到9999999。每条请求用一个协成完成
	//http://player.kuwo.cn/webmusic/st/getNewMuiseByRid?rid=MUSIC_63940553
	//http://antiserver.kuwo.cn/anti.s?rid=MUSIC_65193832&format=aac|mp3&type=convert_url&response=url
	**/

	for i := 2750596; i <= 9999999; i++ {
		go getSongInfo(i)
		time.Sleep(time.Millisecond * 1)

	}
	//getSongInfo(65193832)
}

//遍历数据库，拼接歌曲文件名，下载歌曲到硬盘中
func downloadMusic() error {
	rows, err := db.Query("select * from songinfo")
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer rows.Close()
	var songList []*Song
	for rows.Next() {
		song := new(Song)
		var songId int
		var mp3Url, aacUrl string
		err = rows.Scan(&songId, &song.Music_id, &song.Mv_rid, &song.Name, &song.Song_url, &song.Artist, &song.Artid, &song.Singer, &song.Special, &song.Ridmd591, &song.Mp3size, &song.Artist_url, &song.Auther_url, &song.Playid, &song.Artist_pic, &song.Artist_pic240, &song.Path, &song.Mp3path, &song.Aacpath, &song.Wmadl, &song.Mp3dl, &song.Aacdl, &song.Lyric, &song.Lyric_zz, &song.RequestMp3Url, &song.ReqeustAacUrl, &mp3Url, &aacUrl)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		if song != nil {
			log.Println("requestMp3Url:", song.RequestMp3Url, "requestAacUrl:", song.ReqeustAacUrl)
			songList = append(songList, song)
		} else {
			log.Println("song==nil")
		}

	}
	len := len(songList)
	log.Println("songListLen:", len)

	if len > 0 {
		for _, v := range songList {
			if v != nil {
				var dir = "E:/music/"
				if v.Singer != "" {
					dir = dir + v.Singer
					createFile(dir)
					log.Println("create singer dir:", dir)
				} else {
					log.Println("song.Singer is null")
					return errors.New("song.Singer is null")
				}
				if v.RequestMp3Url != "" {
					go func() {
						mp3UrlRes, err := http.Get(v.RequestMp3Url)
						if err != nil {
							log.Println(err.Error())
							return
						}
						var mp3UrlResBytes []byte
						mp3UrlResBytes, _ = ioutil.ReadAll(mp3UrlRes.Body)
						defer mp3UrlRes.Body.Close()
						mp3Url := string(mp3UrlResBytes)
						log.Println("mp3Url:", mp3Url)
						if mp3Url != "" {
							filePath := dir + "/" + v.Name + "_" + v.Music_id + "_" + v.Special + ".mp3"
							log.Println(filePath)
							downLoad(mp3Url, filePath)
						}

					}()

				}

				if v.ReqeustAacUrl != "" {
					go func() {
						aacUrlRes, err := http.Get(v.ReqeustAacUrl)
						if err != nil {
							log.Println(err.Error())
							return
						}
						var aacUrlResBytes []byte
						aacUrlResBytes, _ = ioutil.ReadAll(aacUrlRes.Body)
						defer aacUrlRes.Body.Close()
						aacUrl := string(aacUrlResBytes)
						log.Println("aacUrl:", aacUrl)
						if aacUrl != "" {
							filePath := dir + "/" + v.Name + "_" + v.Music_id + "_" + v.Special + ".aac"
							log.Println(filePath)
							downLoad(aacUrl, filePath)
						}
					}()
				}

			}

			time.Sleep(time.Millisecond * 200)
		}
	}

	return nil
}

func downLoad(url string, filePath string) error {
	if url != "" && filePath != "" {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		f, err := os.Create(filePath)
		if err != nil {
			return err
		}
		io.Copy(f, res.Body)
		log.Println("下载成功！")
		defer f.Close()
		return nil
	} else {
		return errors.New("url or filePath is illegal")
	}

}

//调用os.MkdirAll递归创建文件夹
func createFile(filePath string) error {
	if !isExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func isExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}