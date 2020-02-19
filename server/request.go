package server

import "gozk/txn"

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  request
 * @Version: 1.0.0
 * @Date: 2020/2/3 4:23 下午
 */

type Request struct {
	SessionId int64
	Cxid      int32
	Type      int32
	Data      []byte
	Protocol  *Protolcol
	TxnHeader *txn.TxnHeader
	Record    interface{}
	//List<Id> authInfo
	Zxid       int64
	CreateTime int64
	Owner      struct{}
}
