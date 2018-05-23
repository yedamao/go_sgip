package errors

import (
	"errors"
)

var (
	SgipLenErr         = errors.New("Operation length error")
	SgipSizeErr        = errors.New("Operation Len larger than MAX_OP_SIZE")
	SgipUserCountErr   = errors.New("Submit UserCount no more than 100")
	SgipBindRespErr    = errors.New("BIND Resp not received")
	SgipMsgPartsLenErr = errors.New("long sms, parts should <= 140")
	SgipBindErr        = errors.New("Bind auth failed.")
	SgipWrongCmdId     = errors.New("resp Wrong CmdId")
)
