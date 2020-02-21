package message

import (
	"encoding/binary"
	"errors"
	"reflect"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  encode
 * @Version: 1.0.0
 * @Date: 2020/2/14 10:36 上午
 */

func EncodePacket(buf []byte, st interface{}) (int, error) {
	v := reflect.ValueOf(st)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return 0, errors.New("st is error")
	}
	return encodePacketValue(buf, v)
}

func encodePacketValue(buf []byte, v reflect.Value) (int, error) {
	//rv := v
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	n := 0
	switch v.Kind() {
	default:
		return 0, errors.New("unknow type")
	case reflect.Bool:
		if v.Bool() {
			buf[n] = 1
		} else {
			buf[n] = 0
		}
		n++
	case reflect.String:
		str := v.String()
		binary.BigEndian.PutUint32(buf[n:n+4], uint32(len(str)))
		copy(buf[n+4:n+4+len(str)], []byte(str))
		n += len(str) + 4
	case reflect.Int32:
		binary.BigEndian.PutUint32(buf[n:n+4], uint32(v.Int()))
		n += 4
	case reflect.Int64:
		binary.BigEndian.PutUint64(buf[n:n+8], uint64(v.Int()))
		n += 8
	case reflect.Struct:
		for i:=0; i<v.NumField(); i++ {
			feild := v.Field(i)
			n2, err := encodePacketValue(buf[n:], feild)
			n += n2
			if err != nil {
				return 0, err
			}
		}
	case reflect.Slice:
		switch v.Type().Elem().Kind() {
		default:
			count := v.Len()
			startN := n
			n += 4
			for i := 0; i < count; i++ {
				n2, err := encodePacketValue(buf[n:], v.Index(i))
				n += n2
				if err != nil {
					return n, err
				}
			}
			binary.BigEndian.PutUint32(buf[startN:startN+4], uint32(count))
		case reflect.Uint8:
			if v.IsNil() {
				binary.BigEndian.PutUint32(buf[n:n+4], uint32(0xffffffff))
				n += 4
			} else {
				bytes := v.Bytes()
				binary.BigEndian.PutUint32(buf[n:n+4], uint32(len(bytes)))
				copy(buf[n+4:n+4+len(bytes)], bytes)
				n += 4 + len(bytes)
			}
		}
	}
	return n, nil
}
