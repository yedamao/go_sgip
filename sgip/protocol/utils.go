package protocol

import (
	"strconv"
	"time"
)

// 节点编号
//
// 通信节点编号规则
// 在整个网关系统中，所有的通信节点(SMG、GNS、SP和SMSC)都有一个唯一的数字编号，不同的SP或SMSC或SMG或GNS编号不能相同，编号由系统管理人员负责分配。编号规则如下：
// SMG的编号规则：1AAAAX
// SMSC的编号规则：	2AAAAX
// SP的编号规则：3AAAAQQQQQ
// GNS的编号规则：4AAAAX
// 其中, AAAA表示四位长途区号(不足四位的长途区号，左对齐，右补零),X表示1位序号,QQQQQ表示5位企业代码。
func NodeId(areaCode, corpId string) (uint32, error) {
	var (
		err error
		ac  int
		ci  int
	)

	// check arg
	if ac, err = strconv.Atoi(areaCode); err != nil {
		return 0, err
	}
	if ci, err = strconv.Atoi(corpId); err != nil {
		return 0, err
	}

	// 0XX 三位区号
	if ac < 100 {
		return uint32(3000000000 + ac*1000000 + ci), nil
	}
	// 0XXX 四位区号
	return uint32(3000000000 + ac*100000 + ci), nil
}

// 格式为十进制的mmddhhmmss，比如11月20日20时32分25秒产生的命令，其第二部分为十进制1120203225
func TimeStamp() uint32 {
	t := time.Now()

	return uint32(int(t.Month())*100000000 + t.Day()*1000000 +
		t.Hour()*10000 + t.Minute()*100 + t.Second())

}
