package message

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  acl
 * @Version: 1.0.0
 * @Date: 2020/2/16 12:42 下午
 */

type ACL struct {
	Perms  int32
	Scheme string
	ID     string
}
