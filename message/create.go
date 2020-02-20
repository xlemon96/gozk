package message

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  create_request
 * @Version: 1.0.0
 * @Date: 2020/2/16 12:42 下午
 */

type CreateRequest struct {
	Path  string
	Data  []byte
	Acl   []*ACL
	Flags int32
}

type CreateResponse struct {
	Path string
}
