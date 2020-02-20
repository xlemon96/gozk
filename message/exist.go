package message

import "gozk/server"

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  exist
 * @Version: 1.0.0
 * @Date: 2020/2/20 5:14 下午
 */

type ExistsRequest struct {
	Path  string
	Watch bool
}

type ExistResponse struct {
	Stat *server.Stat
}
