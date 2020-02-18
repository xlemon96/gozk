package util

import (
	"errors"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  zk_util
 * @Version: 1.0.0
 * @Date: 2020/2/13 3:16 下午
 */

type ZkConn struct {
	conn *zk.Conn
} 

func NewZkConn(hosts []string, timeout time.Duration) (*ZkConn, error) {
	conn, _, err := zk.Connect(hosts, timeout)
	if err != nil {
		return nil, err
	}
	if conn == nil {
		return nil, errors.New("conn is nil")
	}
	zkConn := &ZkConn{conn:conn}
	return zkConn, nil
}

func (this *ZkConn) CreateNode(path string, data []byte, flag int32) (string, error) {
	result, err := this.conn.Create(path, data, flag, zk.WorldACL(zk.PermAll))
	if err != nil {
		return "", err
	}
	return result, err
}

func (this *ZkConn) CreateESNode(path string, acl []zk.ACL) (string, error) {
	result, err := this.conn.CreateProtectedEphemeralSequential(path, []byte(""), acl)
	if err != nil {
		return "", err
	}
	return result, err
}

func (this *ZkConn) GetWithWatcher(path string) (string, <-chan zk.Event, error) {
	result, _, ch, err := this.conn.GetW(path)
	if err != nil {
		return "", nil, err
	}
	return string(result), ch, nil
}

func (this *ZkConn) Exist(path string) (bool, error) {
	exist, _, err := this.conn.Exists(path)
	if err != nil {
		return false, nil
	}
	return exist, nil
}

func (this *ZkConn) Delete(path string) error {
	err := this.conn.Delete(path, -1)
	if err != nil {
		return err
	}
	return nil
}

func (this *ZkConn) GetChildrens(parentPath string) ([]string, error) {
	childrens, _, err := this.conn.Children(parentPath)
	if err != nil {
		return nil, err
	}
	return childrens, nil
}

func (this *ZkConn) Close() {
	this.conn.Close()
}