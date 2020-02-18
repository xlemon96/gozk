package server

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  proto
 * @Version: 1.0.0
 * @Date: 2020/2/4 5:08 下午
 */

type Protolcol struct {
	zookeeperServer *ZookeeperServer
	conn            net.Conn
	reader          *bufio.Reader
	requestChan     chan *InnerRequest
	responseChan    chan struct{}
	Initialized     bool
	Me              struct{}
	SessionId       int64
	SessionTimeout  int32
}

type InnerRequest struct {
	Data      []byte
	Protolcol *Protolcol
}

func NewProtolcol(zookeeperServer *ZookeeperServer, conn net.Conn) *Protolcol {
	protocol := &Protolcol{
		zookeeperServer: zookeeperServer,
		conn:            conn,
		reader:          bufio.NewReader(conn),
		requestChan:     make(chan *InnerRequest),
		Initialized:     false,
		Me:              struct{}{},
	}
	return protocol
}

func (this *Protolcol) Loop() {
	go this.ReadRequestLoop()
	for {
		select {
		case req := <-this.requestChan:
			this.zookeeperServer.ProcessRequest(req)
		}
	}
}

func (this *Protolcol) ReadRequestLoop() {
	for {
		length, _ := this.ReadRequestLength()
		innerRequest, _ := this.ReadRequest(length)
		this.requestChan <- innerRequest
	}
}

//数据包得前四个字节为包长度
func (this *Protolcol) ReadRequestLength() (int32, error) {
	var length int32
	err := binary.Read(this.reader, binary.BigEndian, &length)
	if err != nil {
		return -1, err
	}
	return length, nil
}

func (this *Protolcol) ReadRequest(length int32) (*InnerRequest, error) {
	requestBytes := make([]byte, length)
	n, err := io.ReadFull(this.reader, requestBytes)
	if err != nil {
		return nil, err
	}
	if int32(n) != length {
		//todo,错误待定义，请求体未读完
		return nil, err
	}
	innerRequest := &InnerRequest{
		Data:      requestBytes,
		Protolcol: this,
	}
	return innerRequest, nil
}
