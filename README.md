# WeChatAlarm

> 这是一个用来做微信报警的应用

> 示例
![Alt text](cut.jpeg)


- 首先更改配置文件(config.toml)中的APPID和SECRET然后运行脚本
```sh
./alarm
服务器已启动监听端口为:9200
```

- 第二步 发送post请求
```sh
curl -d"user=123,用户名，多个用户使用逗号分隔&message=测试一下测试一下111" 127.0.0.1:9200

#这时候刚才启动的服务器就会发出请求

POST请求:创建成功 {"access_token":"-gMdaPWaSq4SxZD9sswojqGhGWCOmn7KAtCr3U7kWcHaxXZQ5k2a1j9iJMHcoe10VOu8OAXKfAIAPMZ","expires_in":7200}
当前Token为： 9_bCOmn7KAtCr3U7kWcHaxXZQ5WiDf7X6TkLqiNMg9ubk2a1j9iJMHcoe10VOu8OAXKfAIAPMZ
过期时间还有：7200秒
POST请求:创建成功 {"errcode":0,"errmsg":"ok","msgid":22212}
```
