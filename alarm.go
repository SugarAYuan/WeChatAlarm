package main

import (
	"net/http"
	"strings"
	"fmt"
	"io/ioutil"
	"os"
	"encoding/json"
	"time"
	"github.com/BurntSushi/toml"
)

type AccessToken struct {
	Token string `json:access_token`
	ExpiresIn int `json:expires_in`
	CreateTime int64
}

type Config struct {
	HttpPort string `toml:HttpPort`
	Appid string `toml:"Appid"`
	Secret string `toml:"Secret"`
	TemplateID string `toml:"TemplateID"`
	Openids map[string]string `toml:"Openids"`
}

var accessToken *AccessToken
var config *Config

func main () {

	config = new(Config)
	if _ , err := toml.DecodeFile("./config.toml" , config);err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	http.HandleFunc("/" , alarm)
	fmt.Println("服务器已启动监听端口为:" + config.HttpPort)
	http.ListenAndServe(":"+config.HttpPort , nil)
}

func alarm (w http.ResponseWriter , r *http.Request) {
	r.ParseForm()//解析参数
	sendUserStr := r.PostForm.Get("user")
	message := r.Form.Get("message")
	users := strings.Split(sendUserStr , ",")
	fmt.Println(users , message)

	for _ , u := range users {
		if config.Openids[u] != "" {
			sendTmplete(config.Openids[u] , message)
		}
	}
}
func sendTmplete (openid , sendMsg string) (string , error) {

	reqBody := `
	{
	   "touser":"`+ openid +`",
	   "template_id":"`+config.TemplateID+`",
		  
	   "data":{
			   "first": {
				   "value":"`+sendMsg+`",
				   "color":"#173177"
			   },
			   "remark":{
				   "value":"点击查看详情！",
				   "color":"#173177"
			   }
	   }
	}
`
	token := getAccessToken()
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + token
	res , err  := postReq(url , reqBody)
	return res , err
}

/**
 * @note 发起post请求
 * @params url string 请求url
 * @params reqBody string 请求体
 * @return 请求响应
 */
func postReq (url , reqBody string) (string , error) {
	//创建请求
	req , err := http.NewRequest("POST" , url , strings.NewReader(reqBody))

	if err != nil {
		fmt.Println(err)
		return "" , err
	}

	//增加header
	req.Header.Set("Content-Type", "application/json; encoding=utf-8")

	//执行请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("POST请求:创建请求失败", err)
		return "",err
	}
	//读取响应
	body, err := ioutil.ReadAll(resp.Body) //此处可增加输入过滤
	if err != nil {
		fmt.Println("POST请求:读取body失败", err)
		return "",err
	}

	fmt.Println("POST请求:创建成功", string(body))

	defer resp.Body.Close()

	return string(body),nil
}

/**
 * @note 获取微信的access_token
 * @params *accessToken 用来存放token的结构体
 * @return token string
 */
func  getAccessToken () string {
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + config.Appid + "&secret=" + config.Secret
	if accessToken == nil || (time.Now().Unix() - accessToken.CreateTime) > 7000 {
		accessToken = new(AccessToken)
		res , err := postReq(url , "")

		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		var t map[string]interface{}
		json.Unmarshal([]byte(res) , &t)
		accessToken.Token = fmt.Sprintf("%s" , t["access_token"])
		accessToken.ExpiresIn , _ = fmt.Printf("%d" , t["expires_in"])
		accessToken.CreateTime = time.Now().Unix()
	}

	fmt.Println("当前Token为：" , accessToken.Token)
	fmt.Printf("过期时间还有：%d秒\r\n" , 7200 - (time.Now().Unix() - accessToken.CreateTime))
	return accessToken.Token
}