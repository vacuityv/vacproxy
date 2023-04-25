# vac-go-proxy

用go语言实现的一个http代理工具，支持多端设备



## 功能支持

1、可选的用户名密码验证

2、请求客户端ip白名单设置

3、请求目标域名ip/白名单设置

4、动态更新配置无需重启

## 使用方式

1、下载对应系统的zip包并解压

2、到相应目录下

3、运行

```shell
./vacproxy 
```

会默认使用同一目录下的config.yml配置文件并使用默认的7777端口启动代理

其他启动参数：

```shell
$ vacproxy -help
    Usage of ./vacproxy:
        -bind string
            proxy bind address (default "0.0.0.0:7777")
        -config string
            config file (default "./config.yml")
        -log string
            the log file path (default "./vacproxy.log")
        -pid string
            the pid file path (default "./vacproxy.pid")
        -q  
            quit proxy
        -s string
            Send signal to the daemon:
                stop — shutdown, same as -q
                reload — reloading the configuration file
```

4、停止

```shell
./vacproxy -q
# or
./vacproxy -s -stop
```

## 配置文件说明

```yml
name: test

# 代理鉴权配置，enabled为true且user和password均不为空代表鉴权
auth:
  enabled: true
  user: test
  password: 1234

# 请求ip白名单，放空代表不限制
inAllowList:
#  - 127.0.0.1
#  - 192.168.100.*

# 目标域名/ip白名单，放空代表不限制
outAllowList:
#  - weixin.qq.com
#  - alipay.com
#  - baidu.com
```

更改配置文件后可以使用reload来重新载入，达到免重启生效

```shell
./vacproxy -s reload
```