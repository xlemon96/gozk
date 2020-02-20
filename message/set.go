package message

import "gozk/server"

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  set
 * @Version: 1.0.0
 * @Date: 2020/2/20 5:13 下午
 */

type SetWatchesRequest struct {
	RelativeZxid int64
	DataWatches  []string
	ExistWatches []string
	ChildWatches []string
}

type SetDataResponse struct {
	Stat *server.Stat
}

type SetDAclResponse struct {
	Stat *server.Stat
}
