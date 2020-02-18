package util

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  lock
 * @Version: 1.0.0
 * @Date: 2020/2/13 3:14 下午
 */

type LockInterface interface {
	Lock() error
	UnLock() error
}
