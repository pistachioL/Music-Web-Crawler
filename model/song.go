package model

import "fmt"

//酷狗飙升榜（struct Data 用搜索的结构体）
type Req struct { //原名是Song
	Status int `json:"status"`
	ErrCode int `json:"err_code"`
	Song Song `json:"data"`
	Like bool `json:"like"`
}
type Authors struct {
	AuthorID      string `json:"author_id"`
	AuthorName    string `json:"author_name"`
	IsPublish     string `json:"is_publish"`
	SizableAvatar string `json:"sizable_avatar"`
	Avatar        string `json:"avatar"`
}


//搜索结构体
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
	Song Song `json:"data"`
}
//type Authors struct {
//	AuthorID string `json:"author_id"`
//	AuthorName string `json:"author_name"`
//	IsPublish string `json:"is_publish"`
//	SizableAvatar string `json:"sizable_avatar"`
//	Avatar string `json:"avatar"`
//}
type Song struct { //原名Data
	Id 		 	int 	`db:"id"`
	SongName string 	`json:"song_name"`
	AuthorName string 	`json:"author_name" `
	AlbumName string 	`json:"album_name"`
	PlayURL string 		`json:"play_url"`
	Lyrics string 		`json:"lyrics"`
	Img string 			`json:"img"`
	Timelength int 		`json:"timelength"`
	//Hash string `json:"hash"`

	//Filesize int `json:"filesize"`
	//AudioName string `json:"audio_name"`
	//HaveAlbum int `json:"have_album"`

	//AlbumID string `json:"album_id"`

	//HaveMv int `json:"have_mv"`
	//VideoID string `json:"video_id"`


	//AuthorID string `json:"author_id"`
	//Privilege int `json:"privilege"`
	//Privilege2 string `json:"privilege2"`

	//Authors []Authors `json:"authors"`
	//IsFreePart int `json:"is_free_part"`
	//Bitrate int `json:"bitrate"`
	//RecommendAlbumID string `json:"recommend_album_id"`
	//AudioID string `json:"audio_id"`
	//HasPrivilege bool `json:"has_privilege"`
	//PlayBackupURL string `json:"play_backup_url"`

}

func Save(songName string, authorName string, albumName string, playUrl string, lyrics string,  img string,  timelength int ) {
	db := Conn()
	song := Song{SongName: songName, AuthorName: authorName, AlbumName: albumName,  PlayURL: playUrl,  Lyrics: lyrics, Img: img, Timelength: timelength} // 根据指针找到数据表
	db.SingularTable(true)
	res:= db.Create(&song)
	fmt.Println("歌曲入库结果：", res)
}
