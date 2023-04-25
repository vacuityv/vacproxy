# vac-go-proxy

A http tool implemented by go language, support multi-platform



## Function support

1、optional username and password verification

2、set client ip whitelist/allowlist

3、set target ip/domain whitelist/allowlist

4、reload config without restart

## Usage

1、download the zip package and unzip

2、cd to the directory

3、run

```shell
./vacproxy 
```
this will run the program with default config.yml and 7777 port

Others：

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

4、stop

```shell
./vacproxy -q
# or
./vacproxy -s -stop
```

5、check status

```shell
curl -i --proxy http://test:1234@127.0.0.1:7777 https://bing.com
```

## config.yml

```yml
name: test

# auth config，enabled=true and user and password not "" will verify the credential
auth:
  enabled: true
  user: test
  password: 1234

# client ip whitelist, won't check if none ip here
inAllowList:
#  - 127.0.0.1
#  - 192.168.100.*

# target ip/domain whitelist, won't check if none ip/domain here
outAllowList:
#  - weixin.qq.com
#  - alipay.com
#  - baidu.com
```

you can use reload signal to reload config file without restart program:

```shell
./vacproxy -s reload
```