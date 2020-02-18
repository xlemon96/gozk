package server

import (
	"strings"
	"sync"

	"gozk/message"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  data_tree
 * @Version: 1.0.0
 * @Date: 2020/2/18 9:31 上午
 */

type DataTree struct {
	sync.RWMutex
	Nodes      map[string]*DataNode //所有tree node，key是path："/test/data/test"，value是node指针
	Ephemerals map[int64][]string   //临时节点，key是session id，value是path切片，包含此session创建的所有临时节点
	Root       *DataNode            //根节点
}

func NewDataTree() *DataTree {
	dataTree := &DataTree{
		Nodes:      make(map[string]*DataNode),
		Ephemerals: make(map[int64][]string),
	}
	dataTree.Root = &DataNode{
		Parent:   nil,
		Data:     nil,
		Acl:      -1,
		Children: make([]string, 0),
		Stat:     &Stat{},
	}
	dataTree.Nodes["/"] = dataTree.Root
	dataTree.Nodes[""] = dataTree.Root
	return dataTree
}

func (this *DataTree) GetEphemerals(sessionId int64) []string {
	ephemerals, ok := this.Ephemerals[sessionId]
	if !ok {
		result := make([]string, 0)
		return result
	}
	return ephemerals
}

func (this *DataTree) GetEphemeralsMap() map[int64][]string {
	return this.Ephemerals
}

func (this *DataTree) GetSessions() []int64 {
	sessions := make([]int64, 0)
	for key, _ := range this.Ephemerals {
		sessions = append(sessions, key)
	}
	return sessions
}

func (this *DataTree) AddDataNode(path string, dataNode *DataNode) {
	this.Lock()
	defer this.Unlock()
	this.Nodes[path] = dataNode
}

func (this *DataTree) GetDataNode(path string) *DataNode {
	this.Lock()
	defer this.Unlock()
	return this.Nodes[path]
}

func (this *DataTree) GetNodeCount() int {
	this.Lock()
	defer this.Unlock()
	return len(this.Nodes)
}

func (this *DataTree) GetEphemeralsCount() int {
	var result int
	for _, ephemeral := range this.Ephemerals {
		result += len(ephemeral)
	}
	return result
}

func (this *DataTree) ApproximateDataSize() int64 {
	//todo
	return 0
}

func (this *DataTree) IsSpecialPath(path string) bool {
	if path == "/" {
		return false
	}
	return true
}

func (this *DataTree) CreateNode(path string, data []byte, acl []*message.ACL, ephemeralOwner int64,
	parentCVersion int, zxid int64, time int64) {

	lastSlash := strings.LastIndex(path, "/")
	parentName := path[:lastSlash]
	childName := path[lastSlash+1:]
	stat := &Stat{
		Czxid:          zxid,
		Mzxid:          zxid,
		Ctime:          time,
		Mtime:          time,
		Version:        0,
		Aversion:       0,
		EphemeralOwner: ephemeralOwner,
		Pzxid:          zxid,
	}
	this.Lock()
	defer this.Unlock()
	parentNode, ok := this.Nodes[parentName]
	if !ok {
		//todo,error待定义
		return
	}
	childrens := parentNode.Children
	for _, children := range childrens {
		if children == childName {
			//todo,error待定义
			return
		}
	}
	childrenNode := &DataNode{
		Parent:   parentNode,
		Data:     data,
		Acl:      0,
		Children: make([]string, 0),
		Stat:     stat,
	}
	parentNode.AddChildren(childName)
	this.Nodes[path] = childrenNode
	if ephemeralOwner != 0 {
		ephemeral, ok := this.Ephemerals[ephemeralOwner]
		if !ok {
			ephemeral = make([]string, 0)
			this.Ephemerals[ephemeralOwner] = ephemeral
		}
		this.Ephemerals[ephemeralOwner] = append(this.Ephemerals[ephemeralOwner], path)
	}
	//todo,quota
}

func (this *DataTree) DeleteNode(path string, zxid int64) {
	lastSlash := strings.LastIndex(path, "/")
	parentName := path[:lastSlash]
	childName := path[lastSlash+1:]
	this.Lock()
	defer this.Unlock()
	node, ok := this.Nodes[path]
	if !ok {
		//todo
		return
	}
	delete(this.Nodes, path)
	parentNode, ok := this.Nodes[parentName]
	if !ok {
		//todo
		return
	}
	for index, children := range parentNode.Children {
		if children == childName {
			parentNode.Children = append(parentNode.Children[:index], parentNode.Children[index+1:]...)
			break
		}
	}
	parentNode.Stat.Pzxid = zxid
	eowner := node.Stat.EphemeralOwner
	if eowner != 0 {
		ephemerals := this.Ephemerals[eowner]
		if ephemerals != nil {
			for index, ephemeral := range ephemerals {
				if ephemeral == path {
					this.Ephemerals[eowner] = append(ephemerals[:index], ephemerals[index+1:]...)
					break
				}
			}
		}
	}
	node = nil
}

func (this *DataTree) SetData(path string, data []byte, zxid, time int64, version int32) {
	this.Lock()
	defer this.Unlock()
	node, ok := this.Nodes[path]
	if !ok {
		//todo
		return
	}
	node.Data = data
	node.Stat.Mtime = time
	node.Stat.Mzxid = zxid
	node.Stat.Version = version
}

func (this *DataTree) GetData(path string) []byte {
	this.Lock()
	defer this.Unlock()
	node, ok := this.Nodes[path]
	if !ok {
		//todo
		return nil
	}
	return node.Data
}

func (this *DataTree) GetChildren(path string)[]string {
	this.Lock()
	defer this.Unlock()
	node, ok := this.Nodes[path]
	if !ok {
		//todo
		return nil
	}
	return node.Children
}
