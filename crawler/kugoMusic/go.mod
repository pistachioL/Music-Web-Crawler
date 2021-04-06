module kugomusic

go 1.15

require (
	github.com/astaxie/beego v1.12.3
	github.com/gin-gonic/gin v1.6.3
	github.com/jaydenwen123/go-util v0.0.0-20210115085038-29ef5f0298c0
	github.com/tidwall/gjson v1.7.4
)

require Middlewares v0.0.0

replace Middlewares => ../../middlewares
