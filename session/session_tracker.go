package session

import (
	"fmt"
	"sync"
	"time"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  sessiontracker
 * @Version: 1.0.0
 * @Date: 2020/2/1 11:24 上午
 */

type SessionTracker struct {
	//session map
	Sessions                      map[int64]*Session
	//key is expirationTime, value is  map that key is ptr of Session and value is bool
	SessionWithSameExpirationTime map[int64]map[*Session]bool
	//key id sessionid and value is timeout
	SessionsWithTimeout           map[int64]int32

	//next session id and expirationtime
	nextSessionId      int64
	nextExpirationTime int64

	//expiration interval
	expirationInterval int

	//SessionTracker running flag
	running bool

	//lock
	sync.RWMutex
}

func NewSessionTracker(expirationInterval int, id int64) *SessionTracker {
	//expirationInterval Second convert to Nanosecond (Second = 1000 * 1000 * 1000 Nanosecond)
	interval := expirationInterval * 1000 * 1000 * 1000
	sessionTracker := &SessionTracker{
		Sessions:                      make(map[int64]*Session),
		SessionWithSameExpirationTime: make(map[int64]map[*Session]bool),
		SessionsWithTimeout:           make(map[int64]int32),
		expirationInterval:            interval,
		running:                       true,
	}
	sessionTracker.nextExpirationTime = sessionTracker.RoundToInterval(time.Now().UnixNano())
	sessionTracker.nextSessionId = sessionTracker.InitNextSession(id)
	return sessionTracker
}

//主循环，定期失效session
func (this *SessionTracker) Run() {
	go func() {
		for this.running {
			currentTime := time.Now().UnixNano()
			//未达到失效时间，则休眠到失效时间，失效session
			if this.nextExpirationTime > currentTime {
				sleepTime := time.Duration(this.nextExpirationTime - currentTime)
				time.Sleep(sleepTime * time.Nanosecond)
				continue
			}
			sessionItems, ok := this.SessionWithSameExpirationTime[this.nextExpirationTime]
			if ok {
				delete(this.SessionWithSameExpirationTime, this.nextExpirationTime)
				for sessionItem, _ := range sessionItems {
					fmt.Printf("expire session: %d, expireTime: %d \n", sessionItem.SessionId, sessionItem.ExpireTime)
					this.SetSessionClosing(sessionItem.SessionId)
					//todo, expire the session
				}
			}
			this.nextExpirationTime += int64(this.expirationInterval)
		}
	}()
}

func (this *SessionTracker) CreateSession(timeout int32) int64 {
	this.AddSession(this.nextSessionId, timeout)
	this.nextSessionId++
	return this.nextSessionId
}

func (this *SessionTracker) AddSession(id int64, timeout int32) {
	this.SessionsWithTimeout[id] = timeout
	if _, ok := this.Sessions[id]; !ok {
		session := NewSession(id, 0, timeout)
		this.Sessions[id] = session
	}
	this.TouchSession(id, timeout)
}

func (this *SessionTracker) TouchSession(id int64, timeout int32) bool {
	session, ok := this.Sessions[id]
	if !ok || session.IsClosing {
		return false
	}
	expireTime := this.RoundToInterval(time.Now().UnixNano())
	if session.ExpireTime > expireTime {
		return true
	}
	sessionMap, ok := this.SessionWithSameExpirationTime[session.ExpireTime]
	if ok {
		delete(sessionMap, session)
	}
	session.ExpireTime = expireTime
	sessionMap, ok = this.SessionWithSameExpirationTime[session.ExpireTime]
	if !ok {
		item := make(map[*Session]bool)
		this.SessionWithSameExpirationTime[expireTime] = item
	}
	this.SessionWithSameExpirationTime[session.ExpireTime][session] = true
	return true
}

func (this *SessionTracker) CheckSession(id int64, owner interface{}) bool {
	session, ok := this.Sessions[id]
	if !ok || session.IsClosing {
		return false
	}
	if session.Owner == nil {
		session.Owner = owner
		return true
	}
	if session.Owner != owner {
		return false
	}
	return true
}

func (this *SessionTracker) RemoveSession(id int64) {
	session, ok := this.Sessions[id]
	if ok {
		sessionItems, ok := this.SessionWithSameExpirationTime[id]
		if ok {
			delete(sessionItems, session)
		}
	}
	delete(this.Sessions, id)
	delete(this.SessionsWithTimeout, id)
}

func (this *SessionTracker) SetSessionClosing(id int64) {
	session := this.Sessions[id]
	if session != nil {
		session.IsClosing = false
	}
}

func (this *SessionTracker) SetOwner(id int64, owner interface{}) error {
	session ,ok := this.Sessions[id]
	if !ok {
		//todo, error待定义
		return nil
	}
	session.Owner = owner
	return nil
}

func (this *SessionTracker) RoundToInterval(time int64) int64 {
	tmpExpirationInterval := int64(this.expirationInterval)
	return (time/tmpExpirationInterval + 1) * tmpExpirationInterval
}

func (this *SessionTracker) InitNextSession(id int64) int64 {
	nextSid := 0
	return int64(nextSid)
}

func (this *SessionTracker) ShutDown() {
	this.Lock()
	defer this.Unlock()
	this.running = false
}
