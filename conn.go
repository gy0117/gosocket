package gosocket

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/gy/gosocket/internal/task"
	"github.com/gy/gosocket/internal/xerr"
	"log"
	"net"
	"sync"
	"sync/atomic"
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
	// 当前连接 管理k-v值的
	sm SessionManager
	// 连接是否关闭，0 未关闭，1 关闭
	isClose   int32
	Recycle   func()
	taskQueue *task.TaskQueue
}

// ReadLoop 循环读消息
func (wsConn *WsConn) ReadLoop() {
	wsConn.eventHandler.OnStart(wsConn)
	var err error
	for {
		err = wsConn.readMessage()
		if err != nil {
			break
		}
	}
	wsConn.eventHandler.OnStop(wsConn, err)
	if wsConn.server {
		wsConn.Recycle()
	}
}

func (wsConn *WsConn) GetSessionMap() SessionManager {
	return wsConn.sm
}

// 处理关闭事件
func (wsConn *WsConn) handleCloseEvent(buf *bytes.Buffer) error {
	if atomic.CompareAndSwapInt32(&wsConn.isClose, 0, 1) {
		wsConn.close(buf.String())
	}
	return xerr.NewError(xerr.CloseNormal, nil)
}

func (wsConn *WsConn) handleMessageEvent(msg *Message) error {
	if wsConn.config.OpenUTF8Check && !msg.IsValidText() {
		return xerr.NewError(xerr.ErrCloseUnSupported, errors.New("invalid text encode, must be utf-8 encode"))
	}
	// TODO 消息并行处理
	wsConn.eventHandler.OnMessage(wsConn, msg)
	return nil
}

func (wsConn *WsConn) handleErrorEvent(err error) {
	// 如果conn未关闭，处理error，然后关闭
	if atomic.CompareAndSwapInt32(&wsConn.isClose, 0, 1) {
		ecode := xerr.CloseNormal
		var respErr error
		switch v := err.(type) {
		case *xerr.Error:
			ecode = v.ECode
			respErr = v.Err
		default:
			respErr = err
		}

		content := fmt.Sprintf("ecode: %d, err: %s\n", ecode, respErr.Error())
		wsConn.close(content)
	}
}

func (wsConn *WsConn) close(content string) {
	if err := wsConn.writeMessage(OpcodeConnectionCloseFrame, []byte(content)); err != nil {
		log.Println("conn close and write connection close frame failed, err: ", err)
		return
	}
	if err := wsConn.conn.Close(); err != nil {
		log.Println("conn close failed, err: ", err)
	}
}
