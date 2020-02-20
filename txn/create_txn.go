package txn

import "gozk/message"

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  create_txn
 * @Version: 1.0.0
 * @Date: 2020/2/20 2:02 下午
 */

type CreateTxn struct {
	Path           string
	Data           []byte
	Acl            []*message.ACL
	Ephemeral      bool
	ParentCVersion int32
}
