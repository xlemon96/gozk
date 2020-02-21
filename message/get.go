package message

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  get
 * @Version: 1.0.0
 * @Date: 2020/2/20 5:15 下午
 */

type GetDataRequest struct {
	Path  string
	Watch bool
}

type GetAclRequest struct {
	Path  string
	Watch bool
}

type GetChildreRequest struct {
	Path  string
	Watch bool
}

type GetChildre2Request struct {
	Path  string
	Watch bool
}

type GetDataResponse struct {
	Data []byte
	Stat *Stat
}

type GetAclResponse struct {
	Acl  []*ACL
	Stat *Stat
}

type GetChildrenResponse struct {
	Childrens []string
}

type GetChildren2Response struct {
	Childrens []string
	Stat      *Stat
}
