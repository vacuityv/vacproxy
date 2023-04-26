# vac-go-proxy

[English readme](https://github.com/vacuityv/vac-go-proxy/blob/main/README.en.md)

用go语言实现的一个http代理工具，支持多端设备

## 功能支持

1、可选的用户名密码验证

2、请求客户端ip白名单设置

3、请求目标域名ip/白名单设置

4、动态更新配置无需重启[Windows平台版本不支持]

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
            the pid file path[Windows platform not support] (default "./vacproxy.pid")
        -q  
            quit proxy[Windows platform not support]
        -s string
            Send signal to the daemon[Windows platform not support]:
                stop — shutdown, same as -q
                reload — reloading the configuration file
```

4、停止

windows:

由于目前Windows平台暂不支持后台运行，因此只能 ctrl+c 停止

其他平台:

```shell
./vacproxy -q
# or
./vacproxy -s -stop
```

5、验证运行情况

```shell
# 无鉴权：
curl -i --proxy http://127.0.0.1:7777 https://www.baidu.com
# 有鉴权：
curl -i --proxy http://test:1234@127.0.0.1:7777 https://www.baidu.com
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