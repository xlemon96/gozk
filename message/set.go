package message

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

type SetDataRequest struct {
	Path    string
	Data    []byte
	Version int32
}

type SetAclRequest struct {
	Path    string
	Acl     []ACL
	Version int32
}

type SetDataResponse struct {
	Stat *Stat
}

type SetAclResponse struct {
	Stat *Stat
}
