package server

import "gozk/message"

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  final_request_processor
 * @Version: 1.0.0
 * @Date: 2020/2/18 8:48 下午
 */

type FinalRequestProcessor struct {
	ZookeeperServer *ZookeeperServer
}

func NewFinalRequestProcessor(server *ZookeeperServer) *FinalRequestProcessor {
	finalRequestProcessor := &FinalRequestProcessor{ZookeeperServer:server}
	return finalRequestProcessor
}

func (s *FinalRequestProcessor) ProcessRequest(request *Request)  {
	for {
		if len(s.ZookeeperServer.OutstandingChanges) > 0 &&
			s.ZookeeperServer.OutstandingChanges[0].Zxid <= request.Zxid {
			cr := s.ZookeeperServer.OutstandingChanges[0]
			s.ZookeeperServer.OutstandingChanges = append(s.ZookeeperServer.OutstandingChanges[:0], s.ZookeeperServer.OutstandingChanges[1:]...)
			if cr.Zxid < request.Zxid {
				//todo
			}
			if value, ok := s.ZookeeperServer.OutstandingChangesForPath[cr.Path]; ok {
				if value == cr {
					delete(s.ZookeeperServer.OutstandingChangesForPath, cr.Path)
				}
			}
			continue
		}
		break
	}
	var result *ProcessTxnResult
	if request.TxnHeader != nil {
		result = s.ZookeeperServer.ProcessTxn(request.TxnHeader, request.Record)
	}
	if request.TxnHeader != nil && request.Type == OpCloseSession {
		//todo
		return
	}
	if request.Protocol == nil {
		return
	}
	//String lastOp = "NA";
	//zks.decInProcess();
	//Code err = Code.OK;
	//Record rsp = null;
	//boolean closeSession = false;
	//try {
	//	if (request.hdr != null && request.hdr.getType() == OpCode.error) {
	//	throw KeeperException.create(KeeperException.Code.get((
	//(ErrorTxn) request.txn).getErr()));
	//}
	//
	//	KeeperException ke = request.getException();
	//	if (ke != null && request.type != OpCode.multi) {
	//	throw ke;
	//}
	//
	//	if (LOG.isDebugEnabled()) {
	//	LOG.debug("{}",request);
	//}

	var rsp interface{}
	var err int32
	closeSession := false
	protocol := request.Protocol
	switch request.Type {
	case OpPing:
		replyHeader := &message.ReplyHeader{
			Xid:  -2,
			Zxid: s.ZookeeperServer.LastProcessedZxid(),
			Err:  0,
		}
		protocol.SendResponse(replyHeader, nil)
		return
	case OpCreateSession:
		if err := s.ZookeeperServer.finishSessionInit(protocol); err != nil {
			//todo, print error
		}
		return
	case OpCreate:
		rsp = &message.CreateResponse{Path: ""}
		err = result.Err
	case OpDelete:
		err = result.Err
	case OpSetData:
		rsp = &message.SetDataResponse{Stat:result.Stat}
		err = result.Err
	case OpSetACL:
		rsp = &message.SetDAclResponse{Stat:result.Stat}
		err = result.Err
	case OpCloseSession:
		closeSession = true
		err = result.Err
	case OpSync:
	case OpCheck:
		rsp = &message.SetDataResponse{Stat:result.Stat}
		err = result.Err
	case OpExists:
		existReq := &message.ExistsRequest{}
		_, err :=message.Decode(request.Data, existReq)
		if err != nil {
			//todo
			return
		}
		var stat *Stat
		if existReq.Watch {
			stat = s.ZookeeperServer.DataTree.StatNode(existReq.Path, protocol)
		} else {
			stat = s.ZookeeperServer.DataTree.StatNode(existReq.Path, nil)
		}
		rsp = &message.ExistResponse{Stat:stat}
	case OpGetData:
		getDataReq := &message.GetDataRequest{}
		_, err :=message.Decode(request.Data, getDataReq)
		if err != nil {
			//todo
			return
		}
		node := s.ZookeeperServer.DataTree.GetDataNode(getDataReq.Path)
		if node == nil {
			//todo
		}
		//PrepRequestProcessor.checkACL(zks, zks.getZKDatabase().aclForNode(n),
		//	ZooDefs.Perms.READ,
		//	request.authInfo);
		stat := new(Stat)
		var data []byte
		if getDataReq.Watch {
			data = s.ZookeeperServer.DataTree.GetData(getDataReq.Path, stat, protocol)
		} else {
			data = s.ZookeeperServer.DataTree.GetData(getDataReq.Path, stat, nil)
		}
		rsp = &message.GetDataResponse{
			Data: data,
			Stat: stat,
		}
	case OpGetChildren:
		getChildrenReq := &message.GetChildreRequest{}
		_, err :=message.Decode(request.Data, getChildrenReq)
		if err != nil {
			//todo
			return
		}
		node := s.ZookeeperServer.DataTree.GetDataNode(getChildrenReq.Path)
		if node == nil {
			//todo
		}
		var childres []string
		if getChildrenReq.Watch {
			childres = s.ZookeeperServer.DataTree.GetChildren(getChildrenReq.Path, nil, protocol)
		} else {
			childres = s.ZookeeperServer.DataTree.GetChildren(getChildrenReq.Path, nil, protocol)
		}
		rsp = &message.GetChildrenResponse{
			Childrens: childres,
		}
	case OpGetChildren2:
		getChildren2Req := &message.GetChildreRequest{}
		_, err :=message.Decode(request.Data, getChildren2Req)
		if err != nil {
			//todo
			return
		}
		node := s.ZookeeperServer.DataTree.GetDataNode(getChildren2Req.Path)
		if node == nil {
			//todo
		}
		var childres []string
		if getChildren2Req.Watch {
			childres = s.ZookeeperServer.DataTree.GetChildren(getChildren2Req.Path, nil, protocol)
		} else {
			childres = s.ZookeeperServer.DataTree.GetChildren(getChildren2Req.Path, nil, protocol)
		}
		rsp = &message.GetChildren2Response{
			Childrens: childres,
		}
	case OpGetACL:
	case OpSetWatches:
		setWatchesReq := &message.SetWatchesRequest{}
		_, err :=message.Decode(request.Data, setWatchesReq)
		if err != nil {
			//todo
			return
		}
		s.ZookeeperServer.DataTree.setWatches(setWatchesReq.RelativeZxid,
			setWatchesReq.DataWatches, setWatchesReq.ExistWatches, setWatchesReq.ChildWatches, protocol)
	}
	replyHeader := &message.ReplyHeader{
		Xid:  request.Cxid,
		Zxid: 0,
		Err:  err,
	}
	protocol.SendResponse(replyHeader, rsp)
	if closeSession {
		//todo
	}
}

func (s *FinalRequestProcessor) ShutDown() {
	//todo
}