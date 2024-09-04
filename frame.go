package gosocket

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/gy/gosocket/internal/tools"
	"io"
	"math"
	"unicode/utf8"
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

// 延续帧
type continuationFrame struct {
	// 是否已经初始化
	hasInit  bool
	opcode   Opcode
	buf      *bytes.Buffer
	compress bool
}

type Frame struct {
	Header       [headerFrameLen]byte
	Continuation *continuationFrame
}

func NewFrame() *Frame {
	return &Frame{
		Continuation: &continuationFrame{},
	}
}

// CreateFrame 创建帧数据，帧头 + payload数据
//func (f *Frame) CreateFrame(opcode Opcode, payload []byte) (*bytes.Buffer, error) {
//	//if openUTF8Check && !isValidText(opcode, payload) {
//	//	return nil, internal.NewXError(internal.ErrCloseUnSupported, internal.ErrTextEncode)
//	//}
//	n := len(payload)
//	buf := bufferPool.Get(headerFrameLen + n)
//	header := &Frame{}
//
//}

func (f *Frame) CreateHeader(fin bool, opcode Opcode, server bool, payloadLen int, enableCompress bool) (headerLen int, maskingKey []byte) {
	if fin {
		f.Header[0] |= 0x80
	}

	if enableCompress {
		f.Header[0] |= 0x40
	}

	f.Header[0] |= byte(opcode) & 0x0F

	headerLen = 2
	switch {
	case payloadLen <= 125:
		f.Header[1] |= byte(payloadLen)

	case payloadLen <= math.MaxUint16:
		f.Header[1] |= 126
		binary.BigEndian.PutUint16(f.Header[2:4], uint16(payloadLen))
		headerLen += 2

	default:
		f.Header[1] |= 127
		binary.BigEndian.PutUint64(f.Header[2:10], uint64(payloadLen))
		headerLen += 8
	}

	// 如果需要掩码，则添加掩码键。客户端在发送数据时会随机生成，服务端处理时不需要对数据进行掩码，因为一般为空
	if !server {
		maskingKey, _ = tools.GenerateMaskingKey()
		f.Header[1] |= 128
		binary.LittleEndian.PutUint32(f.Header[headerLen:headerLen+4], binary.LittleEndian.Uint32(maskingKey))
		headerLen += 4
	}
	return
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

// GetFIN 128 或 0，使用 fin != 0 进行判断
func (f *Frame) GetFIN() int {
	return int(f.Header[0] & 0x80)
}

func (f *Frame) GetPayloadLen() int {
	return int(f.Header[1] & 0x7F)
}

func (f *Frame) GetMaskingKey() []byte {
	return f.Header[10:14]
}

func (f *Frame) InitContinuationFrame(opcode Opcode, payloadLen int, compress bool) {
	f.Continuation.hasInit = true
	f.Continuation.opcode = opcode
	f.Continuation.buf = bytes.NewBuffer(make([]byte, 0, payloadLen))
	f.Continuation.compress = compress
}

func (f *Frame) HasInitContinuationFrame() bool {
	return f.Continuation.hasInit
}

func (f *Frame) Write(payloadBytes []byte) {
	f.Continuation.buf.Write(payloadBytes)
}

func (f *Frame) GetContinuationBufLength() int {
	return f.Continuation.buf.Len()
}

func (f *Frame) ResetContinuation() {
	f.Continuation.hasInit = false
	f.Continuation.opcode = 0
	f.Continuation.buf = nil
}

func isValidText(opcode Opcode, data []byte) bool {
	// 连接关闭帧也可能包含可选payload，这个payload中可以包含关闭原因
	if opcode == OpcodeTextFrame || opcode == OpcodeConnectionCloseFrame {
		return utf8.Valid(data)
	}
	return true
}
