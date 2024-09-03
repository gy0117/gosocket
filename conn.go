package gosocket

import (
	"bufio"
	"bytes"
	"github.com/gy/gosocket/internal"
	"net"
	"sync"
)

// WsConn websocket connection
type WsConn struct {
	conn         net.Conn
	bufReader    *bufio.Reader
	eventHandler EventHandler
	frame        *Frame
	config       *Config
	// 标识是否为服务端
	server bool
	lock   sync.Mutex
	sm     SessionManager // 当前连接 管理k-v值的
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

// TODO
func (wsConn *WsConn) handleMessageEvent(msg *Message) error {
	if wsConn.config.OpenUTF8Check && !msg.IsValidText() {
		return internal.NewXError(internal.ErrCloseUnSupported, internal.ErrTextEncode)
	}
	// TODO 消息并行处理

	wsConn.eventHandler.OnMessage(wsConn, msg)
	return nil
}

// TODO
func (wsConn *WsConn) handleErrorEvent(err error) {

}

func (wsConn *WsConn) GetSessionMap() SessionManager {
	return wsConn.sm
}
