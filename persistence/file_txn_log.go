package persistence

import (
	"bufio"
	"os"

	"gozk/message"
	"gozk/txn"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  file_txn_log
 * @Version: 1.0.0
 * @Date: 2020/2/19 1:46 下午
 */

type FileTxnLog struct {
	logFile        *os.File
	logBuf         *logBuf
	streamsToFlush []*logBuf
	lastZxidSeen   int64
	FilePandding   *FilePadding
	Buf            []byte
}

type logBuf struct {
	logBuf *bufio.ReadWriter
	file   *os.File
}

func (p *FileTxnLog) Append(txnHeader *txn.TxnHeader, record interface{}) bool {
	var err error
	if txnHeader == nil {
		return false
	}
	if txnHeader.Zxid <= p.lastZxidSeen {
		//todo
	} else {
		p.lastZxidSeen = txnHeader.Zxid
	}

	if p.logBuf == nil {
		p.logFile, err = os.OpenFile("/Users/jiajianyun/go/src/test/test.txt", os.O_RDWR, 777)
		if err != nil {
			//todo
			return false
		}
		rw := &bufio.ReadWriter{
			Reader: bufio.NewReader(p.logFile),
			Writer: bufio.NewWriter(p.logFile),
		}
		p.logBuf = &logBuf{
			logBuf: rw,
			file:   p.logFile,
		}
		fileHeader := &FileHeader{
			Magic:   0,
			Version: 0,
			DbId:    0,
		}
		n, err := message.EncodePacket(p.Buf[0:], fileHeader)
		if err != nil {
			//todo
			return false
		}
		if _, err := p.logBuf.logBuf.Write(p.Buf[0:n]); err != nil {
			//todo
			return false
		}
		if err := p.logBuf.logBuf.Flush(); err != nil {
			//todo
			return false
		}
		position, _ := p.logFile.Seek(0, 1)
		p.FilePandding.CurrentSize = position
		p.streamsToFlush = append(p.streamsToFlush, p.logBuf)
	}

	p.FilePandding.PadFile(p.logFile)
	//todo
	return true
}

func (p *FileTxnLog) RollLog() error {
	if p.logBuf != nil {
		if err := p.logBuf.logBuf.Flush(); err != nil {
			return err
		}
	}
	p.logBuf = nil
	return nil
}

func (p *FileTxnLog) Commit() error {
	if p.logBuf != nil {
		if err := p.logBuf.logBuf.Flush(); err != nil {
			return err
		}
	}
	for _, buf := range p.streamsToFlush {
		if buf != nil {
			if err := buf.logBuf.Flush(); err != nil {
				return err
			}
		}
	}
	p.streamsToFlush = make([]*logBuf, 0)
	return nil
}

func (p *FileTxnLog) Close() error {
	if p.logBuf != nil {
		if err := p.logBuf.logBuf.Flush(); err != nil {
			return err
		}
		p.logBuf.file.Close()
	}
	for _, buf := range p.streamsToFlush {
		if err := buf.logBuf.Flush(); err != nil {
			return err
		}
		p.logBuf.file.Close()
	}
	return nil
}
