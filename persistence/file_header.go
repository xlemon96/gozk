package persistence

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  file_header
 * @Version: 1.0.0
 * @Date: 2020/2/19 2:03 下午
 */

type FileHeader struct {
	Magic   int32
	Version int32
	DbId    int64
}
