package txn

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  set_data_txn
 * @Version: 1.0.0
 * @Date: 2020/2/20 4:24 下午
 */

type SetDataTxn struct {
	Path    string
	Data    []byte
	Version int32
}
