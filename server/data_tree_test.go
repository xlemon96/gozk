package server

import (
	"fmt"
	"testing"
	"time"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  data_tree_test
 * @Version: 1.0.0
 * @Date: 2020/2/18 10:53 上午
 */

func TestNewDataTree(t *testing.T) {
	tree := NewDataTree()

	tree.CreateNode("/test", []byte("jiajianyun"), nil, 6, 6, 6, 6)

	fmt.Println(string(tree.Nodes["/test"].Data))

	time.Sleep(time.Second * 5)

	tree.DeleteNode("/test", 7)

	node := tree.Nodes["/test"]
	fmt.Println(node)

	time.Sleep(time.Second * 5)
}