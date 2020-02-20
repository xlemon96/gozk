package server

import (
	"math/rand"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  sync_requestProcessor
 * @Version: 1.0.0
 * @Date: 2020/2/19 4:47 下午
 */

type SyncRequestProcessor struct {
	zookeeperServer  *ZookeeperServer
	nextProcessor    ProcessorInterface
	requestsChan     chan *Request
	toFlush          []*Request
	stopChan         chan struct{}
	snapCount        int
	randRoll         int
	snapProcessState int //0未开始，1 已经开始
}

func NewSyncRequestProcessor(zookeeperServer *ZookeeperServer, processor *FinalRequestProcessor) *SyncRequestProcessor {
	syncRequestProcessor := &SyncRequestProcessor{
		zookeeperServer: zookeeperServer,
		requestsChan:    make(chan *Request),
		stopChan:        make(chan struct{}),
		toFlush:         make([]*Request, 0),
		nextProcessor:   processor,
	}
	return syncRequestProcessor
}

func (s *SyncRequestProcessor) Run() {
	go s.loop()
}

func (s *SyncRequestProcessor) loop() {
	logCount := 0
	s.randRoll = rand.Intn(s.snapCount / 2)
	for {
		var request *Request
		select {
		case request = <-s.requestsChan:
			if s.zookeeperServer.FileTxnLog.Append(request.TxnHeader, request.Record) {
				logCount++
				if logCount > (s.randRoll + s.snapCount/2) {
					s.randRoll = rand.Intn(s.snapCount / 2)
					if err := s.zookeeperServer.FileTxnLog.RollLog(); err != nil {
						//todo, print err
						s.stopChan <- struct{}{}
					}
					if s.snapProcessState == 0 {
						go func() {
							s.snapProcessState = 1
							//todo, 快照
							s.snapProcessState = 0
						}()
					}
				}
				//若tuFlush压力不大，则直接处理，否则批量处理
			} else if len(s.toFlush) == 0 {
				s.nextProcessor.ProcessRequest(request)
				continue
			}
			s.toFlush = append(s.toFlush, request)
			if len(s.toFlush) > 1000 {
				s.flush()
			}
		case <-s.stopChan:
			if len(s.toFlush) > 0 {
				s.flush()
			}
			break
		}
	}
}

func (s *SyncRequestProcessor) flush() {
	if len(s.toFlush) == 0 {
		return
	}
	if err := s.zookeeperServer.FileTxnLog.Commit(); err != nil {
		//todo
		return
	}
	for _, req := range s.toFlush {
		if s.nextProcessor != nil {
			s.nextProcessor.ProcessRequest(req)
		}
	}
	s.toFlush = make([]*Request, 0)
}

func (s *SyncRequestProcessor) ProcessRequest(request *Request) {
	s.requestsChan <- request
}

func (s *SyncRequestProcessor) ShutDown() {
	s.stopChan <- struct{}{}
}
