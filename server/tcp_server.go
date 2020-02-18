package server

import (
	"fmt"
	"net"
	"sync"
	"strings"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  tcp_server
 * @Version: 1.0.0
 * @Date: 2020/2/4 5:02 下午
 */

func TCPServer(listener net.Listener, handler TCPHandler) error {
	var wg sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				return fmt.Errorf("listener.Accept() error - %s", err)
			}
			break
		}
		wg.Add(1)
		go func() {
			handler.Handle(conn)
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}