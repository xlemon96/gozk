package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"gozk/txn"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"gozk/message"
	"gozk/persistence"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  zookeeper_server
 * @Version: 1.0.0
 * @Date: 2020/2/3 4:10 下午
 */

type ZookeeperServer struct {
	sync.RWMutex
	ExpirationInterval        int32
	MinSessionTimeout         int32
	MaxSessionTimeout         int32
	SessionTracker            *SessionTracker
	State                     int32
	FirstProcessor            ProcessorInterface
	FileTxnLog                *persistence.FileTxnLog
	OutstandingChanges        []*ChangeRecord
	OutstandingChangesForPath map[string]*ChangeRecord
	DataTree                  *DataTree
	Zxid                      int64
	Protolcols                map[int64]*Protolcol
	//todo
	//private final AtomicLong hzxid = new AtomicLong(0)
	//static final private long superSecret = 0XB3415C00L
	//private ServerCnxnFactory serverCnxnFactory
	//private final ServerStats serverStats
	//private final ZooKeeperServerListener listener
	//private ZooKeeperServerShutdownHandler zkShutdownHandler
}

type ChangeRecord struct {
	Zxid          int64
	Path          string
	Stat          *message.Stat
	ChildrenCount int32
	Acl           []*message.ACL
}

func NewZookeeperServer(minSessionTimeout, maxSessionTimeout, expirationInterval int32) *ZookeeperServer {
	zookeeperServer := &ZookeeperServer{
		ExpirationInterval:        expirationInterval,
		MinSessionTimeout:         minSessionTimeout,
		MaxSessionTimeout:         maxSessionTimeout,
		State:                     ZKINITIAL,
		FirstProcessor:            &PrepRequestProcessor{},
		DataTree:                  NewDataTree(),
		OutstandingChanges:        make([]*ChangeRecord, 0),
		OutstandingChangesForPath: make(map[string]*ChangeRecord),
		Protolcols:                make(map[int64]*Protolcol),
	}
	fileTxnLog := &persistence.FileTxnLog{
		Buf: make([]byte, 65),
		FilePandding: &persistence.FilePadding{
			PreAllocSize: 65535 * 1024,
		},
	}
	zookeeperServer.FileTxnLog = fileTxnLog
	zookeeperServer.CreateSessionTracker()
	return zookeeperServer
}

func (s *ZookeeperServer) Run() error {
	s.StartSessionTracker()
	s.setupRequestProcessors()
	listener, _ := net.Listen("tcp", "localhost:2182")
	handler := &Handler{ZookeeperServer: s}
	err := TCPServer(listener, handler)
	if err != nil {
		return err
	}
	return nil
}

func (s *ZookeeperServer) CreateSessionTracker() {
	sessionTracker := NewSessionTracker(s.ExpirationInterval, 0, s)
	s.SessionTracker = sessionTracker
}

func (s *ZookeeperServer) StartSessionTracker() {
	s.SessionTracker.Run()
}

func (s *ZookeeperServer) ProcessRequest(innerReq *InnerRequest) {
	if innerReq.Protolcol.Initialized {
		s.processRequest(innerReq)
	} else {
		s.processConnectRequest(innerReq)
		innerReq.Protolcol.Initialized = true
	}
}

func (s *ZookeeperServer) processRequest(innerReq *InnerRequest) {
	reqHeader := &message.RequestHeader{}
	n, err := message.Decode(innerReq.Data, reqHeader)
	if err != nil {
		//todo
	} else if n != 8 {
		//todo
	}
	if reqHeader.Type == OpAuth {
		fmt.Printf("get auth packet, %s", innerReq.Protolcol.conn.RemoteAddr())
		authRequest := &message.AuthRequest{}
		_, err := message.Decode(innerReq.Data[n:], authRequest)
		if err != nil {
			//todo
			return
		}
	} else {
		if reqHeader.Type == OpSasl {

		} else {
			req := &Request{
				SessionId: innerReq.Protolcol.SessionId,
				Cxid:      reqHeader.Xid,
				Type:      reqHeader.Type,
				Data:      innerReq.Data[n:],
				Protocol:  innerReq.Protolcol,
				Owner:     innerReq.Protolcol.Me,
				AuthInfo:  innerReq.Protolcol.AuthInfo,
			}
			s.submitRequest(req)
		}
	}
}

func (s *ZookeeperServer) processConnectRequest(innerReq *InnerRequest) error {
	connectReq := &message.ConnectRequest{}
	_, err := message.Decode(innerReq.Data, connectReq)
	if err != nil {
		return err
	} else if connectReq.LastZxidSeen > 0 {
		return errors.New("client lastZxid more than zk lastZxid")
	}
	sessionTimeout := connectReq.TimeOut
	if sessionTimeout < s.MinSessionTimeout {
		sessionTimeout = s.MinSessionTimeout
	} else if sessionTimeout > s.MaxSessionTimeout {
		sessionTimeout = s.MaxSessionTimeout
	}
	innerReq.Protolcol.SessionTimeout = sessionTimeout
	sessionId := connectReq.SessionID
	if sessionId != 0 {
		//todo,客户端发生重连，删除之前的连接，重置session
	}
	s.createSession(innerReq.Protolcol, connectReq.Password, sessionTimeout)
	return nil
}

func (s *ZookeeperServer) createSession(protolcol *Protolcol, password []byte, sessionTimeout int32) {
	timeout := make([]byte, 4)
	binary.BigEndian.PutUint32(timeout, uint32(sessionTimeout))
	sessionId := s.SessionTracker.CreateSession(sessionTimeout)
	protolcol.SessionId = sessionId
	request := &Request{
		SessionId:  sessionId,
		Cxid:       0,
		Type:       OpCreateSession,
		Data:       timeout,
		Protocol:   protolcol,
		CreateTime: time.Now().UnixNano(),
	}
	s.Protolcols[protolcol.SessionId] = protolcol
	s.submitRequest(request)
}

func (s *ZookeeperServer) submitRequest(request *Request) {
	if err := s.touch(request.Protocol); err != nil {
		//todo, drop this request, because session is missing
		return
	}
	s.FirstProcessor.ProcessRequest(request)
}

func (s *ZookeeperServer) GetZxid() int64 {
	return s.Zxid
}

func (s *ZookeeperServer) GetNextZxid() int64 {
	return atomic.AddInt64(&s.Zxid, 1)
}

func (s *ZookeeperServer) setupRequestProcessors() {
	finalP := NewFinalRequestProcessor(s)
	syncP := NewSyncRequestProcessor(s, finalP)
	syncP.Run()
	prepP := NewPrepRequestProcessor(s, syncP)
	prepP.Run()
	s.FirstProcessor = prepP
}

func (s *ZookeeperServer) touch(protolcol *Protolcol) error {
	if protolcol == nil {
		return nil
	}
	sessionId := protolcol.SessionId
	sessionTimeout := protolcol.SessionTimeout
	if !s.SessionTracker.TouchSession(sessionId, sessionTimeout) {
		//todo, error待定义
		return nil
	}
	return nil
}

func (s *ZookeeperServer) finishSessionInit(protolcol *Protolcol) error {
	buf := make([]byte, 256)
	resp := &message.ConnectResponse{
		ProtocolVersion: 0,
		TimeOut:         protolcol.SessionTimeout,
		SessionID:       protolcol.SessionId,
		Password:        nil,
	}
	n, err := message.EncodePacket(buf[4:], resp)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint32(buf[:4], uint32(n))
	_, err = protolcol.conn.Write(buf[:n+4])
	if err != nil {
		return err
	}
	return nil
}

func (s *ZookeeperServer) ProcessTxn(header *txn.TxnHeader, record interface{}) *ProcessTxnResult {
	result := s.DataTree.processTxn(header, record)
	//if (opCode == OpCode.createSession) {
	//	if (txn instanceof CreateSessionTxn) {
	//		CreateSessionTxn cst = (CreateSessionTxn) txn;
	//		sessionTracker.addSession(sessionId, cst
	//		.getTimeOut());
	//	} else {
	//		LOG.warn("*****>>>>> Got "
	//		+ txn.getClass() + " "
	//		+ txn.toString());
	//	}
	//} else if (opCode == OpCode.closeSession) {
	//	sessionTracker.removeSession(sessionId);
	//}
	return result
}

func (s *ZookeeperServer) LastProcessedZxid() int64 {
	return s.DataTree.LastProcessedZxid
}

func (s *ZookeeperServer) Expire(session *Session) {
	sessionId := session.SessionId
	request := &Request{
		SessionId: sessionId,
		Cxid:      0,
		Type:      OpCloseSession,
		Data:      nil,
		Protocol:  nil,
		TxnHeader: nil,
		Record:    nil,
		AuthInfo:  nil,
		Zxid:      0,
	}
	s.submitRequest(request)
}

func (s *ZookeeperServer) removeProtocol(protolcol *Protolcol) {
	s.DataTree.removeProtocol(protolcol)
}
