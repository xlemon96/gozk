package session

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  session
 * @Version: 1.0.0
 * @Date: 2020/2/1 11:09 上午
 */

type Session struct {
	SessionId  int64
	Timeout    int32
	ExpireTime int64
	IsClosing  bool
	Owner      interface{}
}

func NewSession(sessionId, expireTime int64, timeout int32) *Session {
	session := &Session{
		SessionId:  sessionId,
		Timeout:    timeout,
		ExpireTime: expireTime,
		IsClosing:  false,
	}
	return session
}

//func (this *Session) GetSessionId() int64 {
//	return this.SessionId
//}
//
//func (this *Session) GetTimeout() int {
//	return this.Timeout
//}
//
//func (this *Session) IsClose() bool {
//	return this.IsClosing
//}