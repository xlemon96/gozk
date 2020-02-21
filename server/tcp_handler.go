package server

import (
	"fmt"
	"net"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  tcp_handler
 * @Version: 1.0.0
 * @Date: 2020/2/4 5:16 下午
 */

type TCPHandler interface {
	Handle(net.Conn)
}

type Handler struct {
	ZookeeperServer *ZookeeperServer
}

func (s *Handler) Handle(conn net.Conn) {
	fmt.Println("new conn")
	protocol := NewProtolcol(s.ZookeeperServer, conn)
	protocol.Loop()
}
