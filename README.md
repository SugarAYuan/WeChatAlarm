# WeChatAlarm

## alarm 是已经编译好的文件可直接运行

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
#多个user使用,号分隔
curl -d"user=123,zhangsan,xiaowang" 127.0.0.1:9200

#这时候刚才启动的服务器就会发出请求

POST请求:创建成功 {"access_token":"token","expires_in":7200}
当前Token为： token
过期时间还有：7200秒
POST请求:创建成功 {"errcode":0,"errmsg":"ok","msgid":22212}
```
