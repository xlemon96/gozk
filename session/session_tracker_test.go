package session

import (
	"fmt"
	"time"
	"testing"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  sessiontracker_test
 * @Version: 1.0.0
 * @Date: 2020/2/3 11:37 上午
 */

func TestNewSessionTracker(t *testing.T) {
	st := NewSessionTracker(2, 0)

	st.CreateSession(100)
	st.CreateSession(100)

	st.Run()

	time.Sleep(10 * time.Second)

	st.ShutDown()
}

func TestSessionTracker_RemoveSession(t *testing.T) {
	st := NewSessionTracker(2, 0)

	st.CreateSession(100)
	st.CreateSession(100)

	fmt.Println(st.Sessions[0])
	fmt.Println(st.Sessions[1])

	st.RemoveSession(0)

	fmt.Println(st.Sessions[0])
	fmt.Println(st.Sessions[1])
}

