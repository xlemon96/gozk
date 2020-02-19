package txn

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  txn_header
 * @Version: 1.0.0
 * @Date: 2020/2/19 3:23 下午
 */

type TxnHeader struct {
	ClientId int64
	Cxid     int64
	Zxid     int64
	Time     int64
	OpType   int
}
