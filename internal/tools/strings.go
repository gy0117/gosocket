package tools

import (
	"reflect"
	"unsafe"
)

func StringToBytesStandard(s string) []byte {
	return []byte(s)
}

func StringToBytesUnSafe(s string) []byte {
	// 将字符串的地址转换为 *reflect.StringHeader
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))

	// 构造一个 SliceHeader，使用字符串的 Data 指针和长度
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}

	// 将 SliceHeader 转换为 []byte 类型并返回
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func StringToBytesUnSafe2(s string) []byte {
	// strHeader是uintptr数组，模拟字符串头部的布局，strHeader[0]是字符串数据的指针，strHeader[1]是字符串的长度
	strHeader := (*[2]uintptr)(unsafe.Pointer(&s))
	byteSlice := unsafe.Slice((*byte)(unsafe.Pointer(strHeader[0])), strHeader[1])
	return byteSlice
}
