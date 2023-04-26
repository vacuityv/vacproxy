# vac-go-proxy

[English](https://github.com/vacuityv/vac-go-proxy/blob/main/README.en.md)｜[中文](https://github.com/vacuityv/vac-go-proxy/blob/main/README.md)

用go语言实现的一个http代理工具，支持多平台，支持安装为服务启动或者直接启动

## 功能支持

1、可选的用户名密码验证

2、请求客户端ip白名单设置

3、请求目标域名/ip白名单设置


## 使用

使用程序所在目录的 config.yaml 作为配置文件

以下为linux平台命令，其他平台请参照修改

安装为服务：

```shell
./vacproxy install
```

启动服务：

```shell
service vacproxy start
```

停止服务：

```shell
service vacproxy stop
```

重启服务：

```shell
service vacproxy restart
```

卸载服务：

```shell
./vacproxy uninstall
```





验证运行情况

```shell
# 无鉴权：
curl -i --proxy http://127.0.0.1:7777 https://sample.com
# 有鉴权：
curl -i --proxy http://test:1234@127.0.0.1:7777 https://sample.com
```

或者你也可以使用如下参数直接运行：

```shell
./vacproxy -console
```

## 配置文件说明

```yaml
name: vacproxy

bind: 0.0.0.0:7777

# 日志配置/log file config，默认在当前目录的vacproxy.log
#log: /Users/vacuity/log/vacproxy.log

# 代理鉴权配置，enabled为true且user和password均不为空代表鉴权
auth:
  enabled: false
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

更改配置文件后可以使用restart来重启生效