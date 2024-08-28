package gosocket

import (
	"bufio"
	"bytes"
	"github.com/gy/gosocket/internal"
	"net"
)

type Config struct {
	// 读取消息的最大长度
	MaxReadPayloadSize int
	// 是否打开utf-8编码检查
	OpenUTF8Check bool
}

// WsConn websocket connection
type WsConn struct {
	conn         net.Conn
	bufReader    *bufio.Reader
	eventHandler EventHandler
	frame        *Frame
	config       *Config
	// 标识是否为服务端
	isServer bool
}

// ReadLoop 循环读消息
func (wsConn *WsConn) ReadLoop() {
	wsConn.eventHandler.OnStart(wsConn)

	for {
		err := wsConn.readMessage()
		if err != nil {
			break
		}
	}
	// TODO 处理错误
}

// 处理关闭事件
// TODO
func (wsConn *WsConn) handleCloseEvent(buf *bytes.Buffer) error {

	return nil
}

func (wsConn *WsConn) handleMessageEvent(msg *Message) error {
	if wsConn.config.OpenUTF8Check && !msg.IsValidText() {
		return internal.NewXError(internal.ErrCloseUnSupported, internal.ErrTextEncode)
	}
	// TODO 消息并行处理

	wsConn.eventHandler.OnMessage(wsConn, msg)
	return nil
}
