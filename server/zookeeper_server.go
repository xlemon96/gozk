package server

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"

	"gozk/message"
	"gozk/session"
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
	ExpirationInterval int
	MinSessionTimeout  int32
	MaxSessionTimeout  int32
	SessionTracker     *session.SessionTracker
	State              int
	FirstProcessor     ProcessorInterface
	//todo
	//private FileTxnSnapLog txnLogFactory = null
	//private ZKDatabase zkDb
	//private final AtomicLong hzxid = new AtomicLong(0)
	//public final static Exception ok = new Exception("No prob")

	//static final private long superSecret = 0XB3415C00L;
	//private final AtomicInteger requestsInProcess = new AtomicInteger(0);
	//final List<ChangeRecord> outstandingChanges = new ArrayList<ChangeRecord>();
	//// this data structure must be accessed under the outstandingChanges lock
	//final HashMap<String, ChangeRecord> outstandingChangesForPath =
	//new HashMap<String, ChangeRecord>();
	//private ServerCnxnFactory serverCnxnFactory;
	//private final ServerStats serverStats;
	//private final ZooKeeperServerListener listener;
	//private ZooKeeperServerShutdownHandler zkShutdownHandler;

}

func NewZookeeperServer(minSessionTimeout, maxSessionTimeout int32, expirationInterval int) *ZookeeperServer {
	zookeeperServer := &ZookeeperServer{
		ExpirationInterval: expirationInterval,
		MinSessionTimeout:  minSessionTimeout,
		MaxSessionTimeout:  maxSessionTimeout,
		State:              ZKINITIAL,
		FirstProcessor:     &PrepRequestProcessor{},
	}
	zookeeperServer.CreateSessionTracker()
	return zookeeperServer
}

func (this *ZookeeperServer) Run() error {
	this.CreateSessionTracker()
	this.StartSessionTracker()
	listener, _ := net.Listen("tcp", ":2181")
	handler := &Handler{ZookeeperServer: this}
	err := TCPServer(listener, handler)
	if err != nil {
		return err
	}
	return nil
}

func (this *ZookeeperServer) CreateSessionTracker() {
	sessionTracker := session.NewSessionTracker(this.ExpirationInterval, 0)
	this.SessionTracker = sessionTracker
}

func (this *ZookeeperServer) StartSessionTracker() {
	this.SessionTracker.Run()
}

func (this *ZookeeperServer) ProcessRequest(innerReq *InnerRequest) {
	if innerReq.Protolcol.Initialized {
		this.processRequest(innerReq)
	} else {
		this.processConnectRequest(innerReq)
		innerReq.Protolcol.Initialized = true
	}
}

func (this *ZookeeperServer) processRequest(innerReq *InnerRequest) {
	reqHeader := &message.RequestHeader{}
	n, err := message.Decode(innerReq.Data[:7], reqHeader)
	if err != nil {
		//todo
	} else if n != 8 {
		//todo
	}
	if reqHeader.Type == OpAuth {
		fmt.Printf("get auth packet, %s", innerReq.Protolcol.conn.RemoteAddr())
	} else {
		if reqHeader.Type == OpSasl {

		} else {
			req := &Request{
				SessionId: innerReq.Protolcol.SessionId,
				Cxid:      reqHeader.Xid,
				Type:      reqHeader.Type,
				Data:      innerReq.Data,
				Protocol:  innerReq.Protolcol,
				Owner:     innerReq.Protolcol.Me,
			}
			this.submitRequest(req)
		}
	}
}

func (this *ZookeeperServer) processConnectRequest(innerReq *InnerRequest) {
	connectReq := &message.ConnectRequest{}
	_, err := message.Decode(innerReq.Data, connectReq)
	if err != nil {
		//todo
		return
	} else if connectReq.LastZxidSeen > 0 {
		//todo
		return
	}
	sessionTimeout := connectReq.TimeOut
	if sessionTimeout < this.MinSessionTimeout {
		sessionTimeout = this.MinSessionTimeout
	} else if sessionTimeout > this.MaxSessionTimeout {
		sessionTimeout = this.MaxSessionTimeout
	}
	connectReq.TimeOut = sessionTimeout
	sessionId := connectReq.SessionID
	if sessionId != 0 {
		//todo,客户端尝试重置session，server端不允许
	}
	this.createSession(innerReq.Protolcol, connectReq.Password, connectReq.TimeOut)
}

func (this *ZookeeperServer) createSession(protolcol *Protolcol, password []byte, sessionTimeout int32) {
	sessionId := this.SessionTracker.CreateSession(sessionTimeout)
	protolcol.SessionId = sessionId
	request := &Request{
		SessionId:  sessionId,
		Cxid:       0,
		Type:       OpCreateSession,
		Data:       password,
		Protocol:   protolcol,
		CreateTime: time.Now().UnixNano(),
	}
	this.submitRequest(request)
}

func (this *ZookeeperServer) submitRequest(request *Request) {
	this.touch(request.Protocol)
	this.FirstProcessor.ProcessRequest(request)
}

func (this *ZookeeperServer) GetNextZxid() int64 {
	return 0
}

func (this *ZookeeperServer) touch(protolcol *Protolcol) {
	sessionId := protolcol.SessionId
	sessionTimeout := protolcol.SessionTimeout
	if !this.SessionTracker.TouchSession(sessionId, sessionTimeout) {
		//todo, error待定义
		return
	}
}

func (this *ZookeeperServer) finishSessionInit(protolcol *Protolcol) error {
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
