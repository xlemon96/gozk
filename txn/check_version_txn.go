package txn

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  check_version_txn
 * @Version: 1.0.0
 * @Date: 2020/2/20 4:29 下午
 */

type CheckVersionTxn struct {
	Path    string
	Version int32
}
