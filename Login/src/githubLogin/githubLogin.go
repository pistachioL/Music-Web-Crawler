package githubLogin

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Conf struct {
	ClientId 		string
	ClientSecret 	string
	RedirectUrl 	string
}

var conf = Conf {
	ClientId: "cecc9bc83bd8cff1bfb0",
	ClientSecret: "915a9f1814c3c082f23043a19e9b456adc95ece5",
	RedirectUrl: "http://localhost:8080/oauth/redirect",
}


//access_token=179c413581d5e984f020c78b54a3648c469de017&scope=&token_type=bearer
type Token struct {
	AccessToken string `json:"access_token"`
}

func Oauth(w http.ResponseWriter, r *http.Request){
	//解决跨域
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")//header的类型
	w.Header().Set("content-type", "application/json")  //返回数据格式是json
	r.ParseForm()

	var err error
	// 获取授权码
	var code = r.URL.Query().Get("code")

	// 获取 token
	var tokenAuthUrl = GetTokenAuthUrl(code)
	var token *Token
	if token, err = GetToken(tokenAuthUrl); err != nil {
	//	fmt.Println(err)
		return
	}

	//获取用户信息
	var userInfo map[string]interface{}
	if userInfo, err = getUserInfo(token); err != nil {
		fmt.Println("获取用户信息失败", err)
		return
	}

	var userInfoByte []byte
	if userInfoByte, err = json.Marshal(userInfo); err != nil {
		fmt.Println("用户信息转换为[]byte时错误", err)
		return
	}

	//用户信息返回前端
	w.Header().Set("Content-Type","application/json")
	if _, err = w.Write(userInfoByte); err != nil {
		fmt.Println("用户信息返回前端时出错", err)
		return
	}
}



func GetTokenAuthUrl(code string) string {
	return fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		conf.ClientId, conf.ClientSecret, code,
	)
}

// 获取 token
func GetToken(url string) (*Token, error) {
	var request *http.Request
	var err error

	if request, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return nil, err
	}
	request.Header.Set("accept", "application/json")

	//对url发起请求
	var httpClient = http.Client{}
	var response *http.Response
	if response, err = httpClient.Do(request); err != nil {
		return nil, err
	}

	//响应体解析为token
	var token Token
	if err = json.NewDecoder(response.Body).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

func getUserInfo(token *Token) (map[string]interface{}, error) {
	var userInfoUrl = "https://api.github.com/user" //github获取用户信息的api
	var req *http.Request
	var err error
	//初始化请求对象
	if req, err = http.NewRequest(http.MethodGet, userInfoUrl, nil); err != nil {
		return nil, err
	}
	//设置请求头（github的授权token）
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token.AccessToken))

	//发送请求，并获取响应
	var httpClient = http.Client{}
	var res *http.Response
	if res, err = httpClient.Do(req); err != nil {
		return nil, err
	}

	//写入数据 到userInfo
	var userInfo = make(map[string]interface{})
	if err = json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	fmt.Println("userInfo: ", userInfo)
	return userInfo, nil
}

func main() {
	http.HandleFunc("/oauth/redirect", githubLogin.Oauth())
	if err := http.ListenAndServe(":9090", nil); err != nil {
		fmt.Println("监听失败，错误信息为:", err)  // log.Fatal("ListenAndServe: ", err)
		return
	}
}


