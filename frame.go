package gosocket

import (
	"bufio"
	"encoding/binary"
	"io"
)

// Opcode RFC 6455
type Opcode uint8

const (
	OpcodeContinuationFrame Opcode = 0x0 // 延续帧，用于分片消息的后续部分

	OpcodeTextFrame Opcode = 0x1 // 这是一个文本帧，负载数据是文本数据（UTF-8 编码）

	OpcodeBinaryFrame Opcode = 0x2 // 这是一个二进制帧，负载数据是二进制数据

	OpcodeConnectionCloseFrame Opcode = 0x8 // 这是一个连接关闭帧，用于关闭 WebSocket 连接

	OpcodePingFrame Opcode = 0x9 // 这是一个 Ping 帧，用于检测连接的有效性

	OpcodePongFrame Opcode = 0xA // 这是一个 Pong 帧，用于响应 Ping 帧
)

// IsDataFrame 判断操作码是否为数据帧
// 控制帧：0x8到0xF之间；数据帧：0x0到0x2
func (o Opcode) IsDataFrame() bool {
	return o <= OpcodeBinaryFrame
}

// UnMaskPayload 对输入的字节数组payload进行掩码处理
// @0xAAC 亮点
func UnMaskPayload(payload []byte, maskingKey []byte) {
	// 将4字节的掩码键转换为32位无符号整数mk
	var mk = binary.LittleEndian.Uint32(maskingKey)
	// 将mk左移32位，并加上原值，得到64位的ksy64，其前后32位都是相同的掩码键
	// 方便后续对8字节（64位）数据块的处理
	var key64 = uint64(mk)<<32 + uint64(mk)

	// 批量处理64字节块
	for len(payload) >= 64 {
		v := binary.LittleEndian.Uint64(payload) // payload[0:8]
		binary.LittleEndian.PutUint64(payload, v^key64)

		v = binary.LittleEndian.Uint64(payload[8:16])
		binary.LittleEndian.PutUint64(payload[8:16], v^key64)

		v = binary.LittleEndian.Uint64(payload[16:24])
		binary.LittleEndian.PutUint64(payload[16:24], v^key64)

		v = binary.LittleEndian.Uint64(payload[24:32])
		binary.LittleEndian.PutUint64(payload[24:32], v^key64)

		v = binary.LittleEndian.Uint64(payload[32:40])
		binary.LittleEndian.PutUint64(payload[32:40], v^key64)

		v = binary.LittleEndian.Uint64(payload[40:48])
		binary.LittleEndian.PutUint64(payload[40:48], v^key64)

		v = binary.LittleEndian.Uint64(payload[48:56])
		binary.LittleEndian.PutUint64(payload[48:56], v^key64)

		v = binary.LittleEndian.Uint64(payload[56:64])
		binary.LittleEndian.PutUint64(payload[56:64], v^key64)

		// 处理完64字节后，继续处理下一个64字节块
		payload = payload[64:]
	}

	// 剩余字节长度小于64，但大于等于8，批量处理剩余的8字节块
	for len(payload) >= 8 {
		v := binary.LittleEndian.Uint64(payload[:8])
		binary.LittleEndian.PutUint64(payload[:8], v^key64)
		payload = payload[8:]
	}

	var n = len(payload)
	for i := 0; i < n; i++ {
		// 等价于 i % 4， 0000 & 0011，0001 & 0011
		idx := i & 3
		payload[i] ^= maskingKey[idx]
	}
}

// @废弃
// 普通的掩码处理
func applyMask(payload []byte, mask [4]byte) []byte {
	maskedPayload := make([]byte, len(payload))
	for i := range payload {
		maskedPayload[i] = payload[i] ^ mask[i%4]
	}
	return maskedPayload
}

const (
	// 帧头 2～14字节
	headerFrameLen = 14
	// 控制帧的payload的最大字节
	maxControlFramePayloadLen = 125
)

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
	//mask := f.Header[1] & 0x80
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
	if f.GetMask() {
		// 如果mask为1，则读取4字节的masking key，用于解码数据
		if _, err := io.ReadFull(reader, f.Header[10:14]); err != nil {
			return 0, err
		}
	}
	return payloadLen, nil
}

// GetRSV1 RSV1为1，返回true
// 必须先执行ParseHeader方法
func (f *Frame) GetRSV1() bool {
	return f.Header[0]&0x40 != 0
}

func (f *Frame) GetRSV2() bool {
	return f.Header[0]&0x20 != 0
}

func (f *Frame) GetRSV3() bool {
	return f.Header[0]&0x10 != 0
}

// GetMask mask为1，返回true。
// 是否对数据进行掩码处理，客户端发送给服务器的帧必须进行掩码处理，客户端发送的帧应设置为1，服务器发送的帧通常为0
func (f *Frame) GetMask() bool {
	return f.Header[1]&0x80 != 0
}

func (f *Frame) GetOpcode() Opcode {
	return Opcode(f.Header[0] & 0x0F)
}

func (f *Frame) GetFIN() int {
	return int(f.Header[0] & 0x80)
}

func (f *Frame) GetPayloadLen() int {
	return int(f.Header[1] & 0x7F)
}

func (f *Frame) GetMaskingKey() []byte {
	return f.Header[10:14]
}
