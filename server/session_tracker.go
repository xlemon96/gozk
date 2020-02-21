package server

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
	//session map, key is sessionId, value is session of ptr
	Sessions map[int64]*Session
	//key is expirationTime, value is  map that key is ptr of Session and value is bool
	SessionWithSameExpirationTime map[int64]map[*Session]bool
	//key id sessionid and value is timeout
	SessionsWithTimeout map[int64]int32

	//next session id and expirationtime
	nextSessionId      int64
	nextExpirationTime int64

	//expiration interval
	expirationInterval int32

	//SessionTracker running flag
	running bool

	Expire SessionExpirer

	//lock
	sync.RWMutex
}

func NewSessionTracker(expirationInterval int32, id int64, expire SessionExpirer) *SessionTracker {
	//expirationInterval Second convert to Nanosecond (Second = 1000 * 1000 * 1000 Nanosecond)
	interval := expirationInterval * 1000 * 1000 * 1000
	sessionTracker := &SessionTracker{
		Sessions:                      make(map[int64]*Session),
		SessionWithSameExpirationTime: make(map[int64]map[*Session]bool),
		SessionsWithTimeout:           make(map[int64]int32),
		expirationInterval:            interval,
		Expire:                        expire,
		running:                       true,
	}
	sessionTracker.nextExpirationTime = sessionTracker.RoundToInterval(time.Now().UnixNano())
	sessionTracker.nextSessionId = sessionTracker.InitNextSession(id)
	return sessionTracker
}

//主循环，定期失效session
func (s *SessionTracker) Run() {
	go func() {
		for s.running {
			currentTime := time.Now().UnixNano()
			//未达到失效时间，则休眠到失效时间，失效session
			if s.nextExpirationTime > currentTime {
				sleepTime := time.Duration(s.nextExpirationTime - currentTime)
				time.Sleep(sleepTime * time.Nanosecond)
				continue
			}
			sessionItems, ok := s.SessionWithSameExpirationTime[s.nextExpirationTime]
			if ok {
				delete(s.SessionWithSameExpirationTime, s.nextExpirationTime)
				for sessionItem, _ := range sessionItems {
					fmt.Printf("expire session: %d, expireTime: %d \n", sessionItem.SessionId, sessionItem.ExpireTime)
					s.SetSessionClosing(sessionItem.SessionId)
					s.Expire.Expire(sessionItem)
				}
			}
			s.nextExpirationTime += int64(s.expirationInterval)
		}
	}()
}

func (s *SessionTracker) CreateSession(timeout int32) int64 {
	s.AddSession(s.nextSessionId, timeout)
	s.nextSessionId++
	return s.nextSessionId
}

func (s *SessionTracker) AddSession(id int64, timeout int32) {
	s.SessionsWithTimeout[id] = timeout
	if _, ok := s.Sessions[id]; !ok {
		session := NewSession(id, 0, timeout)
		s.Sessions[id] = session
	}
	s.TouchSession(id, timeout)
}

func (s *SessionTracker) TouchSession(id int64, timeout int32) bool {
	session, ok := s.Sessions[id]
	if !ok || session.IsClosing {
		return false
	}
	expireTime := s.RoundToInterval(time.Now().UnixNano())
	if session.ExpireTime > expireTime {
		return true
	}
	sessionMap, ok := s.SessionWithSameExpirationTime[session.ExpireTime]
	if ok {
		delete(sessionMap, session)
	}
	session.ExpireTime = expireTime
	sessionMap, ok = s.SessionWithSameExpirationTime[session.ExpireTime]
	if !ok {
		item := make(map[*Session]bool)
		s.SessionWithSameExpirationTime[expireTime] = item
	}
	s.SessionWithSameExpirationTime[session.ExpireTime][session] = true
	return true
}

func (s *SessionTracker) CheckSession(id int64, owner interface{}) bool {
	session, ok := s.Sessions[id]
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

func (s *SessionTracker) RemoveSession(id int64) {
	session, ok := s.Sessions[id]
	if ok {
		sessionItems, ok := s.SessionWithSameExpirationTime[id]
		if ok {
			delete(sessionItems, session)
		}
	}
	delete(s.Sessions, id)
	delete(s.SessionsWithTimeout, id)
}

func (s *SessionTracker) SetSessionClosing(id int64) {
	session := s.Sessions[id]
	if session != nil {
		session.IsClosing = false
	}
}

func (s *SessionTracker) SetOwner(id int64, owner interface{}) error {
	session, ok := s.Sessions[id]
	if !ok {
		//todo, error待定义
		return nil
	}
	session.Owner = owner
	return nil
}

func (s *SessionTracker) RoundToInterval(time int64) int64 {
	//tmpTime := time * 1000 * 1000 * 1000
	tmpExpirationInterval := int64(s.expirationInterval)
	return (time/tmpExpirationInterval + 1) * tmpExpirationInterval
}

func (s *SessionTracker) InitNextSession(id int64) int64 {
	nextSid := 0
	return int64(nextSid)
}

func (s *SessionTracker) ShutDown() {
	s.Lock()
	defer s.Unlock()
	s.running = false
}
