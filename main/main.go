package main

import (
	"github.com/gin-gonic/gin"
	"githubLogin/crawler"
	"githubLogin/login"
	"githubLogin/middlewares"
	"githubLogin/search"
	"githubLogin/userhome"
)

func start() {
	//todo 程序开始时调用爬虫
}

func main() {
	engine := gin.Default()
	engine.Use(middlewares.Cors())
	engine.Any("/register", login.Register)
	engine.Any("/login", login.Login)
	engine.Any("/oauth/redirect", login.Oauth)
	engine.Any("/editProfile", userhome.UpdateProfile)
	engine.Any("/readProfile", userhome.ReadProfile)
	engine.Any("/popularList", crawler.HandleSongData)
	engine.Any("/search", search.HandleSearch)
	_ = engine.Run(":9091")
}
