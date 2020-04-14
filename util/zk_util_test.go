package util

import (
	"fmt"
	"testing"
	"time"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  zk_util_test
 * @Version: 1.0.0
 * @Date: 2020/2/13 3:24 下午
 */

func TestZkConn_GetChildrens(t *testing.T) {
	hosts := []string{"localhost:2182"}
	zkConn, err := NewZkConn(hosts, time.Second * 5)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(zkConn)
	//fmt.Println(zkConn.GetChildrens("/test"))
	//
	//result.txt, err := zkConn.CreateNode("/test/lock", []byte("lock data"), zk.FlagSequence)
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(zkConn.GetChildrens("/test"))
	//fmt.Println(result.txt)
}

