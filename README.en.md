# vac-go-proxy

[English](https://github.com/vacuityv/vac-go-proxy/blob/main/README.en.md)｜[中文](https://github.com/vacuityv/vac-go-proxy/blob/main/README.md)

A HTTP proxy tool implemented in Go language, supporting multiple platforms, and can be installed as a service or run directly.


## Function support

1. Optional username and password authentication

2. Setting of client IP whitelist for requests

3. Setting of target domain name/IP whitelist for requests

## Usage

Use the config.yaml in the program directory as the configuration file.

The following commands are for the Linux platform. Please modify them for other platforms.

Install as a service:

```shell
./vacproxy install
```

Start the service:

```shell
service vacproxy start
```

Stop the service:

```shell
service vacproxy stop
```

Restart the service:

```shell
service vacproxy restart
```

Uninstall the service:

```shell
./vacproxy uninstall
```

Verify the running status:

```shell
# Without authentication:
curl -i --proxy http://127.0.0.1:7777 https://sample.com
# With authentication:
curl -i --proxy http://test:1234@127.0.0.1:7777 https://sample.com
```

Alternatively, you can use the following parameters to run directly:

```shell
./vacproxy -console
```

## Configuration file

```yaml
bind: 0.0.0.0:7777

# Log configuration/log file configuration, the default file is vacproxy.log in the current directory
#log: /Users/vacuity/log/vacproxy.log

# Proxy authentication configuration, if enabled is true and both user and password are not empty, authentication is enabled
auth:
  enabled: false
  user: test
  password: 1234

# Client IP whitelist for requests, leave empty for no restrictions
inAllowList:
#  - 127.0.0.1
#  - 192.168.100.*

# Target domain name/IP whitelist for requests, leave empty for no restrictions
outAllowList:
#  - weixin.qq.com
#  - alipay.com
#  - baidu.com
```

After changing the configuration file, you can use the restart command to take effect.