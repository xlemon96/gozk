package util

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  lock_test
 * @Version: 1.0.0
 * @Date: 2020/2/13 8:49 下午
 */

var (
	conn *ZkConn
	err  error
)

func init() {
	conn, err = NewZkConn([]string{"localhost:2181"}, time.Second*5)
	if err != nil {
		os.Exit(0)
	}
}

func TestZkLock_Lock(t *testing.T) {
	lock := NewZkLock("/parentlock", conn, zk.WorldACL(zk.PermAll))
	if err := lock.Lock(); err != nil {
		fmt.Println("lock fail")
		os.Exit(0)
	}
	defer func() {
		if err := lock.UnLock(); err != nil {
			fmt.Println(err)
		}
	}()
	fmt.Println("lock success")
	time.Sleep(time.Second * 5)
}