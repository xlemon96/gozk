package main

import (
	"os"

	"gozk/server"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2020/2/18 8:58 下午
 */

func main() {
	zks := server.NewZookeeperServer(5,5,5)
	err := zks.Run()
	if err != nil {
		os.Exit(0)
	}
}