package server

import (
	"strings"
	"time"

	"gozk/message"
	"gozk/txn"
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
		_, err := message.Decode(request.Data, createReq)
		if err != nil {
			return
		}
		s.pRequest2Txn(request, s.zookeeperServer.GetNextZxid(), createReq)
	case OpDelete:
		deleteReq := &message.DeleteRequest{}
		_, err := message.Decode(request.Data, deleteReq)
		if err != nil {
			return
		}
		s.pRequest2Txn(request, s.zookeeperServer.GetNextZxid(), deleteReq)
	case OpSetData:
		setDataReq := &message.SetDataRequest{}
		_, err := message.Decode(request.Data, setDataReq)
		if err != nil {
			return
		}
		s.pRequest2Txn(request, s.zookeeperServer.GetNextZxid(), setDataReq)
	case OpSetACL:
		setAclReq := &message.SetAclRequest{}
		_, err := message.Decode(request.Data, setAclReq)
		if err != nil {
			return
		}
		s.pRequest2Txn(request, s.zookeeperServer.GetNextZxid(), setAclReq)
	case OpCheck:
		checkVersionReq := &message.CheckVersionRequest{}
		_, err := message.Decode(request.Data, checkVersionReq)
		if err != nil {
			return
		}
		s.pRequest2Txn(request, s.zookeeperServer.GetNextZxid(), checkVersionReq)
	case OpMulti:
	case OpCreateSession:
	case OpCloseSession:
		s.pRequest2Txn(request, s.zookeeperServer.GetNextZxid(), nil)
	case OpExists:
	case OpGetData:
	case OpGetChildren:
	case OpGetChildren2:
	case OpGetACL:
	case OpPing:
	case OpSetWatches:
		s.zookeeperServer.SessionTracker.CheckSession(request.SessionId, request.Owner)
	}
	request.Zxid = s.zookeeperServer.GetZxid()
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
		createReq := req.(*message.CreateRequest)
		s.zookeeperServer.SessionTracker.CheckSession(request.SessionId, request.Owner)
		path := createReq.Path
		lastSlash := strings.LastIndex(path, "/")
		if lastSlash == -1 {
			//todo
			return
		}
	case OpCloseSession:
		//todo, 去除临时节点，变更事务队列状态
		s.zookeeperServer.SessionTracker.SetSessionClosing(request.SessionId)
	}
}

func (s *PrepRequestProcessor) ProcessRequest(request *Request) {
	s.requestsChan <- request
}

func (s *PrepRequestProcessor) ShutDown() {
	s.stopChan <- struct{}{}
	s.nextProcessor.ShutDown()
}

//func (s *PrepRequestProcessor) getRecordForPath(path string) *ChangeRecord {
//	lastChange, ok := s.zookeeperServer.OutstandingChangesForPath[path]
//	if !ok {
//		node := s.zookeeperServer.DataTree.GetDataNode(path)
//		if node != nil {
//			children := node.Children
//		}
//	}
//}
