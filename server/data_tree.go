package server

import (
	"strings"
	"sync"

	"gozk/message"
	"gozk/txn"
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
	Nodes             map[string]*DataNode //所有tree node，key是path："/test/data/test"，value是node指针
	Ephemerals        map[int64][]string   //临时节点，key是session id，value是path切片，包含此session创建的所有临时节点
	Root              *DataNode            //根节点
	DataWatches       *WatchManager
	ChildWatches      *WatchManager
	LastProcessedZxid int64
}

type ProcessTxnResult struct {
	ClientId int64
	Cxid     int32
	Zxid     int64
	Err      int32
	Type     int32
	Path     string
	Stat     *Stat
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

func (s *DataTree) GetEphemerals(sessionId int64) []string {
	ephemerals, ok := s.Ephemerals[sessionId]
	if !ok {
		result := make([]string, 0)
		return result
	}
	return ephemerals
}

func (s *DataTree) GetEphemeralsMap() map[int64][]string {
	return s.Ephemerals
}

func (s *DataTree) GetSessions() []int64 {
	sessions := make([]int64, 0)
	for key, _ := range s.Ephemerals {
		sessions = append(sessions, key)
	}
	return sessions
}

func (s *DataTree) AddDataNode(path string, dataNode *DataNode) {
	s.Lock()
	defer s.Unlock()
	s.Nodes[path] = dataNode
}

func (s *DataTree) GetDataNode(path string) *DataNode {
	s.Lock()
	defer s.Unlock()
	return s.Nodes[path]
}

func (s *DataTree) GetNodeCount() int {
	s.Lock()
	defer s.Unlock()
	return len(s.Nodes)
}

func (s *DataTree) GetEphemeralsCount() int {
	var result int
	for _, ephemeral := range s.Ephemerals {
		result += len(ephemeral)
	}
	return result
}

func (s *DataTree) ApproximateDataSize() int64 {
	//todo
	return 0
}

func (s *DataTree) IsSpecialPath(path string) bool {
	if path == "/" {
		return false
	}
	return true
}

func (s *DataTree) CreateNode(path string, data []byte, acl []*message.ACL, ephemeralOwner int64,
	parentCVersion int32, zxid int64, time int64) {

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
	s.Lock()
	defer s.Unlock()
	parentNode, ok := s.Nodes[parentName]
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
	if parentCVersion == -1 {
		parentCVersion = parentNode.Stat.Cversion
		parentCVersion++
	}
	parentNode.Stat.Cversion = parentCVersion //在预处理器已经加一，此处直接赋值即可
	parentNode.Stat.Pzxid = zxid
	//todo，需要处理acl
	childrenNode := &DataNode{
		Parent:   parentNode,
		Data:     data,
		Acl:      0,
		Children: make([]string, 0),
		Stat:     stat,
	}
	parentNode.AddChildren(childName)
	s.Nodes[path] = childrenNode
	if ephemeralOwner != 0 {
		ephemeral, ok := s.Ephemerals[ephemeralOwner]
		if !ok {
			ephemeral = make([]string, 0)
			s.Ephemerals[ephemeralOwner] = ephemeral
		}
		s.Ephemerals[ephemeralOwner] = append(s.Ephemerals[ephemeralOwner], path)
	}
	s.DataWatches.TriggerWatch(path, EventNodeCreated)
	if parentName == "" {
		parentName = "/"
	}
	s.ChildWatches.TriggerWatch(parentName, EventNodeChildrenChanged)
}

func (s *DataTree) DeleteNode(path string, zxid int64) {
	lastSlash := strings.LastIndex(path, "/")
	parentName := path[:lastSlash]
	childName := path[lastSlash+1:]
	s.Lock()
	defer s.Unlock()
	node, ok := s.Nodes[path]
	if !ok {
		//todo
		return
	}
	delete(s.Nodes, path)
	parentNode, ok := s.Nodes[parentName]
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
		ephemerals := s.Ephemerals[eowner]
		if ephemerals != nil {
			for index, ephemeral := range ephemerals {
				if ephemeral == path {
					s.Ephemerals[eowner] = append(ephemerals[:index], ephemerals[index+1:]...)
					break
				}
			}
		}
	}
	node = nil
	if parentName == "" {
		parentName = "/"
	}
	//todo,与原版此处不一样
	s.DataWatches.TriggerWatch(path, EventNodeCreated)
	s.ChildWatches.TriggerWatch(parentName, EventNodeChildrenChanged)
}

func (s *DataTree) SetData(path string, data []byte, zxid, time int64, version int32) {
	s.Lock()
	defer s.Unlock()
	node, ok := s.Nodes[path]
	if !ok {
		//todo
		return
	}
	node.Data = data
	node.Stat.Mtime = time
	node.Stat.Mzxid = zxid
	node.Stat.Version = version
	s.DataWatches.TriggerWatch(path, EventNodeDataChanged)
}

func (s *DataTree) GetData(path string, stat *Stat, protolcol *Protolcol) []byte {
	s.Lock()
	defer s.Unlock()
	node, ok := s.Nodes[path]
	if !ok {
		//todo
		return nil
	}
	//todo,更新node状态
	if protolcol != nil {
		s.DataWatches.AddWatch(path, protolcol)
	}
	return node.Data
}

func (s *DataTree) StatNode(path string, protolcol *Protolcol) *Stat {
	s.Lock()
	defer s.Unlock()
	node, ok := s.Nodes[path]
	if !ok {
		//todo
		return nil
	}
	//todo,更新node状态
	if protolcol != nil {
		s.DataWatches.AddWatch(path, protolcol)
	}
	return node.Stat
}

func (s *DataTree) GetChildren(path string, stat *Stat, protolcol *Protolcol) []string {
	s.Lock()
	defer s.Unlock()
	node, ok := s.Nodes[path]
	if !ok {
		//todo
		return nil
	}
	//todo,跟新配置
	if protolcol != nil {
		s.ChildWatches.AddWatch(path, protolcol)
	}
	return node.Children
}

func (s *DataTree) processTxn(header *txn.TxnHeader, record interface{}) *ProcessTxnResult {
	result := &ProcessTxnResult{
		ClientId: header.ClientId,
		Cxid:     header.Cxid,
		Zxid:     header.Zxid,
		Err:      0,
		Type:     header.Type,
	}
	switch header.Type {
	case OpCreate:
		createTxn := record.(*txn.CreateTxn)
		result.Path = createTxn.Path
		var ephemeralOwner int64
		if !createTxn.Ephemeral {
			ephemeralOwner = result.ClientId
		} else {
			ephemeralOwner = 0
		}
		s.CreateNode(createTxn.Path, createTxn.Data, createTxn.Acl, ephemeralOwner, createTxn.ParentCVersion,
			header.Zxid, header.Time)
	case OpDelete:
		deleteTxn := record.(*txn.DeleteTxn)
		result.Path = deleteTxn.Path
		s.DeleteNode(deleteTxn.Path, header.Zxid)
	case OpSetData:
		setDataTxn := record.(*txn.SetDataTxn)
		result.Path = setDataTxn.Path
		s.SetData(setDataTxn.Path, setDataTxn.Data, header.Time, header.Zxid, setDataTxn.Version)
	case OpSetACL:
	case OpCloseSession:
		s.killSession(header.ClientId, header.Zxid)
	case OpError:
		errTxn := record.(*txn.ErrorTxn)
		result.Err = errTxn.Err
	case OpCheck:
		checkTxn := record.(*txn.CheckVersionTxn)
		result.Path = checkTxn.Path
	case OpMulti:
		//todo
	}
	if result.Zxid > s.LastProcessedZxid {
		s.LastProcessedZxid = result.Zxid
	}
	//todo
	//if (header.getType() == OpCode.create &&
	//	rc.err == Code.NODEEXISTS.intValue()) {
	//	LOG.debug("Adjusting parent cversion for Txn: " + header.getType() +
	//		" path:" + rc.path + " err: " + rc.err);
	//	int lastSlash = rc.path.lastIndexOf('/');
	//	String parentName = rc.path.substring(0, lastSlash);
	//	CreateTxn cTxn = (CreateTxn)txn;
	//	try {
	//		setCversionPzxid(parentName, cTxn.getParentCVersion(),
	//			header.getZxid());
	//	} catch (KeeperException.NoNodeException e) {
	//		LOG.error("Failed to set parent cversion for: " +
	//			parentName, e);
	//		rc.err = e.code().intValue();
	//	}
	//} else if (rc.err != Code.OK.intValue()) {
	//	LOG.debug("Ignoring processTxn failure hdr: " + header.getType() +
	//		" : error: " + rc.err);
	//}
	return result
}

func (s *DataTree) killSession(sessionId, zxid int64) {
	paths, ok := s.Ephemerals[sessionId]
	if ok {
		delete(s.Ephemerals, sessionId)
		for _, path := range paths {
			s.DeleteNode(path, zxid)
		}
	}
}

func (s *DataTree) removeProtocol(protolcol *Protolcol) {
	s.DataWatches.RemoveWatch(protolcol)
	s.ChildWatches.RemoveWatch(protolcol)
}

func (s *DataTree) clear() {
	s.Root = nil
	s.Nodes = make(map[string]*DataNode)
	s.Ephemerals = make(map[int64][]string)
}

func (s *DataTree) setWatches(relativeZxid int64, dataWatches, existWatches, childWatches []string, protolcol *Protolcol) {
	for _, path := range dataWatches {
		node,ok := s.Nodes[path]
		if !ok {
			event := &WatcherEvent{
				Type:  EventNodeDeleted,
				State: StateConnected,
				Path:  path,
			}
			protolcol.Process(event)
		} else if node.Stat.Mzxid > relativeZxid {
			event := &WatcherEvent{
				Type:  EventNodeDataChanged,
				State: StateConnected,
				Path:  path,
			}
			protolcol.Process(event)
		} else {
			s.DataWatches.AddWatch(path, protolcol)
		}
	}

	for _, path := range existWatches {
		_,ok := s.Nodes[path]
		if ok {
			event := &WatcherEvent{
				Type:  EventNodeCreated,
				State: StateConnected,
				Path:  path,
			}
			protolcol.Process(event)
		} else {
			s.DataWatches.AddWatch(path, protolcol)
		}
	}

	for _, path := range childWatches {
		node,ok := s.Nodes[path]
		if !ok {
			event := &WatcherEvent{
				Type:  EventNodeDeleted,
				State: StateConnected,
				Path:  path,
			}
			protolcol.Process(event)
		} else if node.Stat.Mzxid > relativeZxid {
			event := &WatcherEvent{
				Type:  EventNodeChildrenChanged,
				State: StateConnected,
				Path:  path,
			}
			protolcol.Process(event)
		} else {
			s.ChildWatches.AddWatch(path, protolcol)
		}
	}
}