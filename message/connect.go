package message

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  connect
 * @Version: 1.0.0
 * @Date: 2020/2/18 5:21 下午
 */

type ConnectRequest struct {
	ProtocolVersion int32
	LastZxidSeen    int64
	TimeOut         int32
	SessionID       int64
	Password        []byte
}

type ConnectResponse struct {
	ProtocolVersion int32
	TimeOut         int32
	SessionID       int64
	Password        []byte
}
