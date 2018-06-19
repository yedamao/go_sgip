# SGIP 1.2
[![Build Status](https://travis-ci.org/yedamao/go_sgip.svg?branch=master)](https://travis-ci.org/yedamao/go_sgip)
[![Go Report Card](https://goreportcard.com/badge/github.com/yedamao/go_sgip)](https://goreportcard.com/report/github.com/yedamao/go_sgip)
[![codecov](https://codecov.io/gh/yedamao/go_sgip/branch/master/graph/badge.svg)](https://codecov.io/gh/yedamao/go_sgip)

go_sgip是为SP设计实现的SGIP 1.2协议开发工具包。包括sgip协议包和命令行工具。

## 安装
```
go get github.com/yedamao/go_sgip/...
cd $GOPATH/src/github.com/yedamao/go_sgip && make
```

## Sgip协议包

###  support operation

- [x] Bind
- [x] BindResp
- [x] Unbind
- [x] UnbindResp
- [x] Submit
- [x] SubmitResp
- [x] Deliver
- [x] DeliverResp
- [x] Report
- [x] ReportResp

## transmitter

## receiver
 
### Example
参照cmd/transmitter/main.go, cmd/receiver/main.go

## 命令行工具

### mockserver
SMG短信网关模拟器

```
Usage of ./bin/mockserver:
  -addr string
        监听地址 (default ":8801")
```

### transmitter
使用短链接提交短信

```
Usage of ./bin/transmitter:
  -area-code string
        长途区号 (default "010")
  -corp-id string
        5位企业代码 (default "00000")
  -dest-number string
        接收手机号码, 86..., 多个使用，分割
  -host string
        SMSC host (default "localhost")
  -msg string
        短信内容
  -name string
        Login Name
  -passwd string
        Login Password
  -port int
        SMSC port (default 8801)
  -service-type string
        业务代码，由SP定义
  -sp-number string
        SP的接入号码
```

### receiver
负责接收运营商上行短信及状态消息

```
Usage of ./bin/receiver:
  -addr string
        上行监听地址 (default ":8001")
  -count int
        worker 数量 (default 5)
```


