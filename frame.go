package gows

import (
	"bufio"
	"encoding/binary"
	"io"
)

// 帧头 2～14字节
const headerFrameLen = 14

type Frame struct {
	Header [headerFrameLen]byte
}

// ParseHeader 解析帧头，获取payload len
// https://developer.mozilla.org/zh-CN/docs/Web/API/WebSockets_API/Writing_WebSocket_servers
func (f *Frame) ParseHeader(reader *bufio.Reader) (int, error) {
	// 读取前两个字节
	if _, err := io.ReadFull(reader, f.Header[0:2]); err != nil {
		return 0, err
	}

	//fin := f.Header[0] & 0x80
	// 0x0F 二进制是0000 1111
	//opcode := f.Header[0] & 0x0F
	mask := f.Header[1] & 0x80
	// 0111 1111
	payloadLen := int(f.Header[1] & 0x7F)

	// 读取扩展的 payload len
	if payloadLen == 126 {
		if _, err := io.ReadFull(reader, f.Header[2:4]); err != nil {
			return 0, err
		}
		payloadLen = int(binary.BigEndian.Uint16(f.Header[2:4]))

	} else if payloadLen == 127 {
		if _, err := io.ReadFull(reader, f.Header[2:10]); err != nil {
			return 0, err
		}
		payloadLen = int(binary.BigEndian.Uint64(f.Header[2:10]))

	}

	// 是否对数据进行掩码处理，客户端发送给服务器的帧必须进行掩码处理
	if mask == 1 {
		// 如果mask为1，则读取4字节的masking key，用于解码数据
		if _, err := io.ReadFull(reader, f.Header[10:14]); err != nil {
			return 0, err
		}
	}
	return payloadLen, nil
}
