package server

import (
	"bufio"
	"io"
	"net"
	"encoding/binary"

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
		}
	}
}

func (s *Protolcol) ReadRequestLoop() {
	for {
		length, _ := s.ReadRequestLength()
		innerRequest, _ := s.ReadRequest(length)
		s.requestChan <- innerRequest
	}
}

//数据包得前四个字节为包长度
func (s *Protolcol) ReadRequestLength() (int32, error) {
	var length int32
	err := binary.Read(s.reader, binary.BigEndian, &length)
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
	n1, err := message.EncodePacket(buf[4:], rh)
	if err != nil {
		return
	}
	n2, err := message.EncodePacket(buf[4+n1:], cr)
	if err != nil {
		return
	}
	binary.BigEndian.PutUint32(buf[0:3], uint32(n1+n2))
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
