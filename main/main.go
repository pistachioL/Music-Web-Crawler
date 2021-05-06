package main

import (
	"github.com/gin-gonic/gin"
	"githubLogin/collect"
	"githubLogin/crawler"
	"githubLogin/login"
	"githubLogin/middlewares"
	"githubLogin/redis"
	"githubLogin/search"
	"githubLogin/userhome"
)

func start() {
	//todo 程序开始时调用爬虫
}

func main() {
	engine := gin.Default()
	engine.Use(middlewares.Cors())
	//search.SaveSearchKey()
	//crawler.GetSongDetail()
	engine.Any("/register", login.Register)
	engine.Any("/login", login.Login)
	engine.Any("/oauth/redirect", login.Oauth)
	engine.Any("/editProfile", userhome.UpdateProfile)
	engine.Any("/readProfile", userhome.ReadProfile)

	engine.Any("/popularList", crawler.HandleSongData)
	engine.Any("/popularListDB", crawler.HandleDBdata)
	engine.Any("/kuwoPopularList", crawler.KuwoSongList)
	engine.Any("/qqPopularList", crawler.QQSongList)


	engine.Any("/setRecentlyPlay", redis.SetRecentPlay)
	engine.Any("/getRecentlyPlay", redis.GetRecentPlay)
	engine.Any("/search", search.HandleSearch)

	engine.Any("/addCollection", collect.AddCollection)
	engine.Any("/queryCollection", collect.QueryCollection)
	engine.Any("/queryFiveCollection", collect.QueryFiveCollection)
	engine.Any("/cancelCollection", collect.CancelCollection)

	_ = engine.Run(":9091")
}
