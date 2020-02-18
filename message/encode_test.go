package message

import (
	"fmt"
	"testing"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  encode_test
 * @Version: 1.0.0
 * @Date: 2020/2/14 10:53 上午
 */

type Student struct {
	Name    string `json:"name"`
	Age     int32  `json:"age"`
	Address string `json:"address"`
}

func TestEncodePacket(t *testing.T) {
	s := &Student{
		Name:    "jiajianyun",
		Age:     27,
		Address: "test",
	}

	buf := make([]byte, 256)

	n, err := EncodePacket(buf, s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(n)

	n, err = Decode(buf[:n], &Student{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(n)
}
