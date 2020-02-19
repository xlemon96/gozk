package persistence

import "os"

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  file_padding
 * @Version: 1.0.0
 * @Date: 2020/2/19 2:25 下午
 */

type FilePadding struct {
	PreAllocSize int64
	CurrentSize  int64
}

func (p *FilePadding) PadFile(file *os.File) {
	position, err := file.Seek(0,1)
	if err != nil {
		//todo
		return
	}
	newFileSize := p.calculateFileSizeWithPadding(position, p.CurrentSize, p.PreAllocSize)
	if newFileSize != p.CurrentSize {
		//todo,填充文件
		//_, err := file.WriteAt([]byte("0"), newFileSize)
		//if err != nil {
		//	//todo
		//	return
		//}
		p.CurrentSize = newFileSize
	}
}

/*
	CurrentSize表示为当前文件已填充大小
	计算是否需要扩容
	若需要，则扩充为2 * PreAllocSize
*/
func (p *FilePadding) calculateFileSizeWithPadding(position, filesize, preAllocSize int64) int64 {
	if (preAllocSize >0 && position + 4096 >= filesize) {
		if (position > filesize) {
			filesize = position + preAllocSize
			filesize -= filesize % preAllocSize
		} else {
			filesize += preAllocSize
		}
	}
	return filesize
}