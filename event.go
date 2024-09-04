package gosocket

import (
	"bytes"
	"unicode/utf8"
)

type Message struct {
	// 消息的操作码
	Opcode Opcode
	// 消息内容
	Content *bytes.Buffer
	// 消息是否压缩
	compress bool
}

// IsValidText 如果opcode为0x1，则payload必须时utf-b编码的文本数据（规定）
func (m *Message) IsValidText() bool {
	// 连接关闭帧也可能包含可选payload，这个payload中可以包含关闭原因
	if m.Opcode == OpcodeTextFrame || m.Opcode == OpcodeConnectionCloseFrame {
		return utf8.Valid(m.Content.Bytes())
	}
	return true
}

type EventHandler interface {
	// OnStart 建立连接事件
	OnStart(conn *WsConn)
	// OnPing ping
	OnPing(conn *WsConn, payload []byte)
	// OnPong pong
	OnPong(conn *WsConn, payload []byte)
	// OnMessage 发送消息
	OnMessage(conn *WsConn, msg *Message)
	// OnStop 关闭连接
	OnStop(conn *WsConn, err error)
}
