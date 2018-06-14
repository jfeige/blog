package models

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math"
	rd "math/rand"
	"reflect"
	"strings"
	"time"
)

/**
检测 指定的value值是否存在Array，Slice，或者Map的value中
*/

func InArray(obj interface{}, target interface{}) bool {

	target_tp := reflect.TypeOf(target)
	target_vl := reflect.ValueOf(target)

	switch target_tp.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < target_vl.Len(); i++ {
			if obj == target_vl.Index(i).Interface() {
				return true
			}
		}
	case reflect.Map:
		for _, v := range target_vl.MapKeys() {
			if obj == target_vl.MapIndex(v).Interface() {
				return true
			}
		}
	}

	return false
}

/**
随机数(all,char,number)
*/
func RandStr(length int, format ...string) string {
	var tp = "all"
	if len(format) > 0 && format[0] != "" {
		tp = strings.ToLower(format[0])
	}
	var bytes []byte
	var r *rd.Rand
	var result []byte
	switch tp {
	case "char":
		bytes = []byte("abcdefghijklmnopqrstuvwxyz")
		if length > len(bytes) {
			length = len(bytes)
		}
		r = rd.New(rd.NewSource(time.Now().UnixNano()))
		for i := 0; i < length; i++ {
			result = append(result, bytes[r.Intn(len(bytes))])
		}
	case "number":
		bytes = []byte("0123456789")
		if length > len(bytes) {
			length = len(bytes)
		}
		r = rd.New(rd.NewSource(time.Now().UnixNano()))
		for i := 0; i < length; i++ {
			result = append(result, bytes[r.Intn(len(bytes))])
		}
	default:
		bytes = []byte("abcdefghijklmnopqrstuvwxyz0123456789")
		if length > len(bytes) {
			length = len(bytes)
		}
		r = rd.New(rd.NewSource(time.Now().UnixNano()))
		for i := 0; i < length; i++ {
			result = append(result, bytes[r.Intn(len(bytes))])
		}
	}

	return string(result)
}

/**
md5
*/
func ToMd5(obj interface{}) string {
	var ret string
	switch obj := obj.(type) {
	case string:
		data := []byte(obj)
		ret = fmt.Sprintf("%x", md5.Sum(data))
	case []byte:
		ret = fmt.Sprintf("%x", md5.Sum(obj))
	case int:
		t := int32(obj)
		buffer := bytes.NewBuffer([]byte{})
		binary.Write(buffer, binary.BigEndian, t)
		ret = fmt.Sprintf("%x", md5.Sum(buffer.Bytes()))
	case int64:
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(obj))
		ret = fmt.Sprintf("%x", md5.Sum(buffer))
	case int32:
		buffer := make([]byte, 4)
		binary.BigEndian.PutUint32(buffer, uint32(obj))
		ret = fmt.Sprintf("%x", md5.Sum(buffer))
	case float64:
		buffer := make([]byte, 8)
		bits := math.Float64bits(obj)
		binary.LittleEndian.PutUint64(buffer, bits)
		ret = fmt.Sprintf("%x", md5.Sum(buffer))
	case float32:
		buffer := make([]byte, 4)
		bits := math.Float32bits(obj)
		binary.LittleEndian.PutUint32(buffer, bits)
		ret = fmt.Sprintf("%x", md5.Sum(buffer))
	default:
	}
	return ret
}
