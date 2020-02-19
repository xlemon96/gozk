package server

import (
	"gozk/message"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  prep_request_processor
 * @Version: 1.0.0
 * @Date: 2020/2/3 4:22 下午
 */

type PrepRequestProcessor struct {
	//zk server
	zookeeperServer *ZookeeperServer

	nextProcessor ProcessorInterface

	//请求channel
	requestsChan chan *Request

	//loop停止channel
	stopChan chan struct{}
}

func NewPrepRequestProcessor(zookeeperServer *ZookeeperServer) *PrepRequestProcessor {
	prepRequestProcessor := &PrepRequestProcessor{
		zookeeperServer: zookeeperServer,
	}
	return prepRequestProcessor
}

func (this *PrepRequestProcessor) Run() {
	go this.loop()
}

func (this *PrepRequestProcessor) loop() {
	for {
		select {
		case req := <-this.requestsChan:
			this.pRequest(req)
		case <-this.stopChan:
			close(this.requestsChan)
			break
		}
	}
}

func (this *PrepRequestProcessor) pRequest(request *Request) {
	switch request.Type {
	case OpCreate:
		createReq := &message.CreateRequest{}
		_, _ = message.Decode(request.Data, createReq)
		this.pRequest2Txn(request, this.zookeeperServer.GetNextZxid(), createReq)
	case OpCreateSession:
	case OpCloseSession:

	}
	request.Zxid = this.zookeeperServer.GetNextZxid()
	this.nextProcessor.ProcessRequest(request)
}

func (this *PrepRequestProcessor) pRequest2Txn(request *Request, zxid int64, req interface{})  {
	switch request.Type {
	case OpCreate:
		req = req.(*message.CreateRequest)
		this.zookeeperServer.SessionTracker.CheckSession(request.SessionId, request.Owner)
	}
}

func (this *PrepRequestProcessor) ProcessRequest(request *Request) {
	this.requestsChan <- request
}

func (this *PrepRequestProcessor) ShutDown() {
	this.stopChan <- struct{}{}
}
