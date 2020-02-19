package persistence

import (
	"fmt"
	"os"
	"testing"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  file_padding_test
 * @Version: 1.0.0
 * @Date: 2020/2/19 2:51 下午
 */

func TestFilePadding_PadFile(t *testing.T) {
	file, _ := os.OpenFile("/Users/jiajianyun/go/src/test/test.txt", os.O_RDWR, 777)
	p := &FilePadding{
		PreAllocSize: 10,
		CurrentSize:  3,
	}
	p.PadFile(file)
}

func TestFilePadding_calculateFileSizeWithPadding(t *testing.T) {

	p := &FilePadding{}
	fmt.Println(p.calculateFileSizeWithPadding(5555, 5552, 4096))

}