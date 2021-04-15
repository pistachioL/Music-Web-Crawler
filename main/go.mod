module main

go 1.15

require github.com/gin-gonic/gin v1.6.3

require Login v0.0.0

replace Login => ../login/githubLogin

require Middlewares v0.0.0

replace Middlewares => ../middlewares

require (
	Crawler v0.0.0
	github.com/gogf/gf v1.15.6
)

replace Crawler => ../crawler/kugoMusic
