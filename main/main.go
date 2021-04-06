package main

import (
	"Crawler"
	"Middlewares"
	"github.com/gin-gonic/gin"
)


func main() {
	//第三方登录
	//http.HandleFunc("/oauth/redirect", githubLogin.Oauth)
	//if err := http.ListenAndServe(":9090", nil); err != nil {
	//	fmt.Println("监听失败，错误信息为:", err)  // log.Fatal("ListenAndServe: ", err)
	//	return
	//}

	engine := gin.Default()
	engine.Use(middlewares.Cors())
	engine.Any("/popularList", kugoMusic.HandleSongData)
	engine.Run(":9091")
}