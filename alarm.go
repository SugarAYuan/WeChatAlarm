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
	Token string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
	CreateTime int64
}

type Config struct {
	HttpPort string `toml:"HttpPort"`
	Appid string `toml:"Appid"`
	Secret string `toml:"Secret"`
	TemplateIDs map[string]string `toml:"TemplateIDs"`
	Openids map[string]string `toml:"Openids"`
}

type sendTemplateData struct {
	Touser string `json:"touser"`
	TemplateId string `json:"template_id"`
	Url string `json:"url"`
	Data map[string]map[string]string `json:"data"`
}

type requestData struct {
	User string `json:"user"`
	Message map[string]string `json:"message"`
	TemplateKey string `json:"template_key"`
	Url string `json:"url"`
}

var accessToken *AccessToken
var config *Config
var templateResult map[string]map[string]map[string]string
func main () {

	config = new(Config)
	if _ , err := toml.DecodeFile("./config.toml" , config);err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//获取所有模板
	templateData , err := ioutil.ReadFile("./template.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	json.Unmarshal([]byte(templateData) , &templateResult)

	http.HandleFunc("/" , alarm)
	fmt.Println("服务器已启动监听端口为:" + config.HttpPort)
	http.ListenAndServe(":"+config.HttpPort , nil)
}

func alarm (w http.ResponseWriter , r *http.Request) {
	//解析发送过来的json数据
	requestData := new(requestData)
	err := json.NewDecoder(r.Body).Decode(requestData)
	users := strings.Split(requestData.User , ",")

	if err != nil {
		fmt.Println(err)
		return
	}

	//获取模板ID
	templateID := config.TemplateIDs[requestData.TemplateKey]
	//解析发送过来的message
	for k , v := range requestData.Message {
		templateResult[templateID][k]["value"] = v
	}

	re := make(map[string]string)
	//查找用户发送报警
	for _ , u := range users {
		if config.Openids[u] != "" {
			//发送
			e , _ := sendTmplete(config.Openids[u] , templateID , requestData.Url)
			//发送结果记录
			re[u] = e
		}
	}
	b , err := json.Marshal(re)
	if err != nil {
		fmt.Fprintf(w , "error")
	}
	//返回发送结果
	fmt.Fprintf(w , string(b))
}
func sendTmplete (openid , templateID , backUrl string) (string , error) {
	//解析模板消息
	sendTemplateData := new(sendTemplateData)
	sendTemplateData.Touser = openid
	sendTemplateData.TemplateId = templateID
	sendTemplateData.Url = backUrl
	sendTemplateData.Data = templateResult[templateID]
	reqBody , err := json.Marshal(sendTemplateData)
	if err != nil {
		fmt.Println(err)
		return "" , err
	}
	fmt.Println(string(reqBody))
	token := getAccessToken()
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + token
	res , err  := postReq(url , string(reqBody))
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