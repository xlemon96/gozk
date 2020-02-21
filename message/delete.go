package message

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  delete
 * @Version: 1.0.0
 * @Date: 2020/2/20 9:27 下午
 */

type DeleteRequest struct {
	Path    string
	Version int32
}

type DeleteResponse struct{

}