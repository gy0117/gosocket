package gosocket

import (
	"bufio"
	"bytes"
	"net"
)

type Config struct {
	// 读取消息的最大长度
	MaxReadPayloadSize int
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
func (wsConn *WsConn) handleClose(buf *bytes.Buffer) error {
	return nil
}
