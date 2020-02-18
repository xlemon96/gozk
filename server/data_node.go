package server

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  data_node
 * @Version: 1.0.0
 * @Date: 2020/2/18 9:26 上午
 */

type DataNode struct {
	Parent   *DataNode
	Data     []byte
	Acl      int64
	Children []string
	Stat     *Stat
}

func (this *DataNode) AddChildren(child string) {
	this.Children = append(this.Children, child)
}
