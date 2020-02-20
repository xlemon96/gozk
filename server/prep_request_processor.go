package server

import (
	"gozk/message"
	"gozk/txn"
	"time"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  prep_request_processor
 * @Version: 1.0.0
 * @Date: 2020/2/3 4:22 下午
 */

type PrepRequestProcessor struct {
	//zkserver
	zookeeperServer *ZookeeperServer

	nextProcessor ProcessorInterface

	//请求channel
	requestsChan chan *Request

	//loop停止channel
	stopChan chan struct{}
}

func NewPrepRequestProcessor(zookeeperServer *ZookeeperServer, processor *SyncRequestProcessor) *PrepRequestProcessor {
	prepRequestProcessor := &PrepRequestProcessor{
		zookeeperServer: zookeeperServer,
		requestsChan:    make(chan *Request),
		stopChan:        make(chan struct{}),
	}
	prepRequestProcessor.nextProcessor = processor
	return prepRequestProcessor
}

func (s *PrepRequestProcessor) Run() {
	go s.loop()
}

func (s *PrepRequestProcessor) loop() {
	for {
		select {
		case req := <-s.requestsChan:
			s.pRequest(req)
		case <-s.stopChan:
			break
		}
	}
}

func (s *PrepRequestProcessor) pRequest(request *Request) {
	switch request.Type {
	case OpCreate:
		createReq := &message.CreateRequest{}
		_, _ = message.Decode(request.Data, createReq)
		s.pRequest2Txn(request, s.zookeeperServer.GetNextZxid(), createReq)
	case OpCreateSession:
	case OpCloseSession:

	}
	request.Zxid = s.zookeeperServer.GetNextZxid()
	s.nextProcessor.ProcessRequest(request)
}

func (s *PrepRequestProcessor) pRequest2Txn(request *Request, zxid int64, req interface{}) {
	request.TxnHeader = &txn.TxnHeader{
		ClientId: request.SessionId,
		Cxid:     request.Cxid,
		Zxid:     request.Zxid,
		Type:     request.Type,
		Time:     time.Now().UnixNano(),
	}
	switch request.Type {
	case OpCreate:
		req = req.(*message.CreateRequest)
		s.zookeeperServer.SessionTracker.CheckSession(request.SessionId, request.Owner)
	}
}

func (s *PrepRequestProcessor) ProcessRequest(request *Request) {
	s.requestsChan <- request
}

func (s *PrepRequestProcessor) ShutDown() {
	s.stopChan <- struct{}{}
}
