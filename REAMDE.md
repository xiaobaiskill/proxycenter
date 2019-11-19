goProxyCenter
====

### 简介
```
    通过ip代理池获取ip放入workpool，通过workpool 请求的程序 
```

### 快速启动
`docker-compose up -d`

### 访问方式
```
curl -X POST http://127.0.0.1:3000/index?proxy=true -d '{"action":"https://jimqaweb.mlytics.ai/cache.txt"}'
```


### refer:
[參考專案go代理池](https://github.com/henson/proxypool)


### 备注
```
mysql ：
    port: 3307
    user: root
    passwd: mysqlroot
    database : proxycenter
```