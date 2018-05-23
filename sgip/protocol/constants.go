package protocol

const (
	SGIP_VERSION = 0x12

	// 系统中每个消息包最大不超过2K字节
	MAX_OP_SIZE = 2048

	// 群发接收手机号码最大为100
	MAX_USER_COUNT
)

// Command ID
const (
	SGIP_BIND = 0x00000001 + iota
	SGIP_UNBIND
	SGIP_SUBMIT
	SGIP_DELIVER
	SGIP_REPORT
	SGIP_ADDSP
	SGIP_MODIFYSP
	SGIP_DELETESP
	SGIP_QUERYROUTE
	SGIP_ADDTELESEG
	SGIP_MODIFYTELESEG
	SGIP_DELETETELESEG
	SGIP_ADDSMG
	SGIP_MODIFYSMG
	SGIP_DELETESMG
	SGIP_CHECKUSER
	SGIP_USERRPT
)

const (
	SGIP_BIND_REP = 0x80000001 + iota
	SGIP_UNBIND_REP
	SGIP_SUBMIT_REP
	SGIP_DELIVER_REP
	SGIP_REPORT_REP
	SGIP_ADDSP_REP
	SGIP_MODIFYSP_REP
	SGIP_DELETESP_REP
	SGIP_QUERYROUTE_REP
	SGIP_ADDTELESEG_REP
	SGIP_MODIFYTELESEG_REP
	SGIP_DELETETELESEG_REP
	SGIP_ADDSMG_REP
	SGIP_MODIFYSMG_REP
	SGIP_DELETESMG_REP
	SGIP_CHECKUSER_REP
	SGIP_USERRPT_REP
)

const (
	STAT_OK RespStatus = iota // 无错误，命令正确接收

	// 1-20所指错误一般在各类命令的应答中用到

	STAT_ILLLOGIN  // 非法登录，如登录名、口令出错、登录名与口令不符等
	STAT_RPTLOGIN  // 重复登录，如在同一TCP/IP连接中连续两次以上请求登录
	STAT_MUCHCONN  // 连接过多，指单个节点要求同时建立的连接数过多
	STAT_ERLGNTYPE // 登录类型错，指bind命令中的logintype字段出错
	STAT_ERARGFMT  // 参数格式错，指命令中参数值与参数类型不符或与协议规定的范围不符
	STAT_ILLUSRNUM // 非法手机号码，协议中所有手机号码字段出现非86130号码或手机号码前未加“86”时都应报错
	STAT_ERSEQ     // 消息ID错
	STAT_ERLEN     // 信息长度错
	STAT_ILLSEQ    // 非法序列号，包括序列号重复、序列号格式错误等
	STAT_ILLOPGNS  // 非法操作GNS
	STAT_NODEBUSY  // 节点忙，指本节点存储队列满或其他原因，暂时不能提供服务的情况

	// 21-32所指错误一般在report命令中用到

	STAT_DSTCNTRCH = 21 + iota // 目的地址不可达，指路由表存在路由且消息路由正确但被路由的节点暂时不能提供服务的情况
	STAT_ROUTER                // 路由错，指路由表存在路由但消息路由出错的情况，如转错SMG等
	STAT_ROUTENEST             // 路由不存在，指消息路由的节点在路由表中不存在
	STAT_INVCHGNUM             // 计费号码无效，鉴权不成功时反馈的错误信息
	STAT_USRCNTRCH             // 用户不能通信（如不在服务区、未开机等情况）
	STAT_MEMFULL               // 手机内存不足
	STAT_NTSPTSMS              // 手机不支持短消息
	STAT_RCVERR                // 手机接收短消息出现错误
	STAT_UNKNUSR               // 不知道的用户
	STAT_NTSPTFUN              // 不提供此功能
	STAT_ILLDEV                // 非法设备
	STAT_SYSFAIL               // 系统失败
	STAT_SMSCFULL              // 短信中心队列满
)

// MessageCoding
const (
	ASCII = 0  // 纯ASCII字符串
	UCS2  = 8  // UCS2编码
	GBK   = 15 // GBK编码
)
