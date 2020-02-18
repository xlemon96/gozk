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
	ZKINITIAL  = 1
	ZKRUNNING  = 2
	ZKSHUTDOWN = 3
	ZKERROR    = 4
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
