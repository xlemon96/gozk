package message

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  auth
 * @Version: 1.0.0
 * @Date: 2020/2/19 8:49 下午
 */

type AuthRequest struct {
	Type   int32
	Scheme string
	Auth   []byte
}
