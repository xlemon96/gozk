package server

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