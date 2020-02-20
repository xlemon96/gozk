package server

import (
	"gozk/message"
	"gozk/txn"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  request
 * @Version: 1.0.0
 * @Date: 2020/2/3 4:23 下午
 */

type Request struct {
	SessionId  int64
	Cxid       int32 //客户端请求头id，表示不同请求之间的顺序
	Type       int32
	Data       []byte
	Protocol   *Protolcol
	TxnHeader  *txn.TxnHeader //事务头
	Record     interface{}
	AuthInfo   []*message.ID //鉴权信息
	Zxid       int64         //事务id
	CreateTime int64
	Owner      struct{}
}
