package util

import (
	"strconv"
	"errors"
	"strings"
	"fmt"

	"github.com/samuel/go-zookeeper/zk"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  zk_lock
 * @Version: 1.0.0
 * @Date: 2020/2/14 10:20 上午
 */

var (
	ErrDeadlock  = errors.New("zk: trying to acquire a lock twice")
	ErrNoNode    = errors.New("zk: node does not exist")
)

type ZkLock struct {
	conn     *ZkConn
	path     string
	acl      []zk.ACL
	lockPath string
	seq      int
}

func NewZkLock(path string, conn *ZkConn, acl []zk.ACL) *ZkLock {
	return &ZkLock{
		conn: conn,
		path: path,
		acl:  acl,
	}
}

func parseSeq(path string) (int, error) {
	items := strings.Split(path, "-")
	return strconv.Atoi(items[len(items)-1])
}

func (this *ZkLock) Lock() error {
	if this.lockPath != "" {
		return ErrDeadlock
	}
	prefix := fmt.Sprintf("%s/lock-", this.path)
	path := ""
	var err error
	for i := 0; i < 3; i++ {
		path, err = this.conn.CreateESNode(prefix, this.acl)
		if err == nil {
			break
		} else if err.Error() == ErrNoNode.Error() {
			items := strings.Split(this.path, "/")
			tmpPath := ""
			for _, item := range items[1:] {
				tmpPath = "/" + item
				ok, err := this.conn.Exist(tmpPath)
				if err != nil {
					return err
				}
				if ok {
					continue
				}
				_, err = this.conn.CreateNode(tmpPath, []byte{}, 0)
				if err != nil {
					return err
				}
			}
		} else {
			return err
		}
	}
	if err != nil {
		return err
	}
	seq, err := parseSeq(path)
	if err != nil {
		return err
	}
	for {
		childrens, err := this.conn.GetChildrens(this.path)
		if err != nil {
			return err
		}
		lowestSeq := seq
		prevSeq := -1
		prevSeqPath := ""
		for _, children := range childrens {
			s, err := parseSeq(children)
			if err != nil {
				return err
			}
			if s < lowestSeq {
				lowestSeq = s
			}
			if s < seq && s > prevSeq {
				prevSeq = s
				prevSeqPath = children
			}
		}
		if seq == lowestSeq {
			break
		}
		_, ch, err := this.conn.GetWithWatcher(this.path + "/" + prevSeqPath)
		if err != nil && err != ErrNoNode {
			return err
		} else if err != nil && err == ErrNoNode {
			continue
		}
		ev := <- ch
		if ev.Err != nil {
			return ev.Err
		}
	}
	this.seq = seq
	this.lockPath = path
	return nil
}

func (this *ZkLock) UnLock() error {
	if err := this.conn.Delete(this.path); err != nil {
		return err
	}
	this.lockPath = ""
	this.seq = 0
	return nil
}

