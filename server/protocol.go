package server

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"

	"gozk/message"
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
	stop            chan struct{}
	lenBytes        []byte
	Initialized     bool
	Me              struct{}
	SessionId       int64
	SessionTimeout  int32
	AuthInfo        []*message.ID
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
		lenBytes:        make([]byte, 4),
		stop:            make(chan struct{}),
		Initialized:     false,
		Me:              struct{}{},
	}
	return protocol
}

func (s *Protolcol) Loop() {
	go s.ReadRequestLoop()
	for {
		select {
		case req := <-s.requestChan:
			s.zookeeperServer.ProcessRequest(req)
		case <-s.stop:
			if err := s.conn.Close(); err != nil {
				//todo, print error
			}
			break
		}
	}
}

func (s *Protolcol) ReadRequestLoop() {
	var err error
	var length int32
	var innerRequest *InnerRequest
	for {
		length, err = s.ReadRequestLength()
		if err != nil {
			goto ERR
		}
		innerRequest, err = s.ReadRequest(length)
		if err != nil {
			goto ERR
		}
		s.requestChan <- innerRequest
	}
ERR:
	if err != io.EOF {
		//todo, print log, eof是由于客户端关闭了连接，server端返回eof异常
	}
	s.stop <- struct{}{}
}

//数据包得前四个字节为包长度
func (s *Protolcol) ReadRequestLength() (int32, error) {
	_, err := io.ReadFull(s.reader, s.lenBytes)
	length := int32(binary.BigEndian.Uint32(s.lenBytes[:4]))
	if err != nil {
		return -1, err
	}
	return length, nil
}

func (s *Protolcol) ReadRequest(length int32) (*InnerRequest, error) {
	requestBytes := make([]byte, length)
	n, err := io.ReadFull(s.reader, requestBytes)
	if err != nil {
		return nil, err
	}
	if int32(n) != length {
		//todo,错误待定义，请求体未读完
		return nil, err
	}
	innerRequest := &InnerRequest{
		Data:      requestBytes,
		Protolcol: s,
	}
	return innerRequest, nil
}

func (s *Protolcol) SendResponse(rh *message.ReplyHeader, cr interface{}) {
	buf := make([]byte, 256)
	var n1, n2 int
	var err error
	n1, err = message.EncodePacket(buf[4:], rh)
	if err != nil {
		return
	}
	if cr != nil {
		n2, err = message.EncodePacket(buf[4+n1:], cr)
		if err != nil {
			return
		}
	}
	binary.BigEndian.PutUint32(buf[:4], uint32(n1+n2))
	_, err = s.conn.Write(buf[:n1+n2+4])
	if err != nil {
		return
	}
}

func (s *Protolcol) Process(event *WatcherEvent) {
	replyHeader := &message.ReplyHeader{
		Xid:  -1,
		Zxid: -1,
		Err:  0,
	}
	e := &Event{
		Type:  int32(event.Type),
		State: int32(event.State),
		Path:  event.Path,
	}
	s.SendResponse(replyHeader, e)
}
