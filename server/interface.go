package server

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  interface
 * @Version: 1.0.0
 * @Date: 2020/2/5 12:02 下午
 */

type ProcessorInterface interface {
	ProcessRequest(request *Request)
	ShutDown()
}

type SessionExpirer interface {
	Expire(session *Session)
}
