package server

import (
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

func (this *Handler) Handle(conn net.Conn) {
	protocol := NewProtolcol(this.ZookeeperServer, conn)
	protocol.Loop()
}
