package server

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  const
 * @Version: 1.0.0
 * @Date: 2020/2/4 3:25 下午
 */

//zookeerper server state
const (
	ZKINITIAL  int32 = 1
	ZKRUNNING  int32 = 2
	ZKSHUTDOWN int32 = 3
	ZKERROR    int32 = 4
)

//op code
const (
	OpNotification  = 0
	OpCreate        = 1
	OpDelete        = 2
	OpExists        = 3
	OpGetData       = 4
	OpSetData       = 5
	OpGetACL        = 6
	OpSetACL        = 7
	OpGetChildren   = 8
	OpSync          = 9
	OpPing          = 11
	OpGetChildren2  = 12
	OpCheck         = 13
	OpMulti         = 14
	OpAuth          = 100
	OpSetWatches    = 101
	OpSasl          = 102
	OpCreateSession = -10
	OpCloseSession  = -11
	OpError         = -1
)

type EventType int32

const (
	EventNodeCreated         EventType = 1
	EventNodeDeleted         EventType = 2
	EventNodeDataChanged     EventType = 3
	EventNodeChildrenChanged EventType = 4
	EventSession             EventType = -1
	EventNotWatching         EventType = -2
)

var (
	eventNames = map[EventType]string{
		EventNodeCreated:         "EventNodeCreated",
		EventNodeDeleted:         "EventNodeDeleted",
		EventNodeDataChanged:     "EventNodeDataChanged",
		EventNodeChildrenChanged: "EventNodeChildrenChanged",
		EventSession:             "EventSession",
		EventNotWatching:         "EventNotWatching",
	}
)

type State int32

const (
	StateUnknown           State = -1
	StateDisconnected      State = 0
	StateConnecting        State = 1
	StateAuthFailed        State = 4
	StateConnectedReadOnly State = 5
	StateSaslAuthenticated State = 6
	StateExpired           State = -112
	StateConnected         State = 100
	StateHasSession        State = 101
)
