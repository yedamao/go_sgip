# SGIP 1.2
[![Build Status](https://travis-ci.org/yedamao/go_sgip.svg?branch=master)](https://travis-ci.org/yedamao/go_sgip)

This is an implementation of SGIP 1.2 for Go

## protocol
decode/encode sgip operation

### support operation

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



## cmd

### transmitter

使用短链接提交短信

usage:

```
./transmitter -host=localhost -port=8801 -area-code=010 -corp-id=10690 -service-type=your-type -sp-number=106900000001 -msg=hahatest -name=xxxxxx -passwd=xxxxxx -dest-number=8617600000000
```


### receiver

负责接收运营商上行短信及状态消息

usage:

```
./receiver -addr=":8001" -count=5
```
