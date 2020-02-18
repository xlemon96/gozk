package message

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  request_header_test
 * @Version: 1.0.0
 * @Date: 2020/2/7 5:55 下午
 */

var data = &RequestHeader{
	Xid:  1,
	Type: 1,
}

func TestRequestHeader_Encode(t *testing.T) {
	d, _ := data.Encode()
	fmt.Println(len(d))
	reader := bytes.NewReader(d)
	buf := bufio.NewReaderSize(reader, 30)
	tt := make([]byte, 60)
	n, _ := buf.Read(tt)
	fmt.Println(len(tt))
	fmt.Println(n)
}

func TestRequestHeader_Decode(t *testing.T) {
	d, _ := data.Encode()
	decode_data := &RequestHeader{}
	_ = decode_data.Decode(d)
	fmt.Println(decode_data)
}
